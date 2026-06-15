package subscription

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shubham071122/collab/internal/config"
)

type Service struct {
	subRepo *Repository
	cfg     *config.Config
}

func NewService(subRepo *Repository, cfg *config.Config) *Service {
	return &Service{
		subRepo: subRepo,
		cfg:     cfg,
	}
}

type RazorpaySubscriptionRequest struct {
	PlanID         string            `json:"plan_id"`
	TotalCount     int               `json:"total_count"`
	Quantity       int               `json:"quantity"`
	CustomerNotify int               `json:"customer_notify"`
	Notes          map[string]string `json:"notes"`
}

type RazorpaySubscriptionResponse struct {
	ID       string `json:"id"`
	ShortURL string `json:"short_url"`
	Status   string `json:"status"`
}

type RazorpayWebhookPayload struct {
	Event   string `json:"event"`
	Payload struct {
		Subscription struct {
			Entity struct {
				ID         string            `json:"id"`
				Status     string            `json:"status"`
				CurrentEnd *int64            `json:"current_end"`
				Notes      map[string]string `json:"notes"`
			} `json:"entity"`
		} `json:"subscription"`
		Payment struct {
			Entity struct {
				ID             string `json:"id"`
				Amount         int    `json:"amount"`
				Status         string `json:"status"`
				SubscriptionID string `json:"subscription_id"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
}

func (s *Service) GetSubscriptionByUserID(userID string) (*UserSubscription, error) {
	return s.subRepo.GetSubscriptionByUserID(userID)
}

func (s *Service) CreateDefaultSubscription(userID string) error {
	return s.subRepo.CreateDefaultSubscription(userID)
}

func (s *Service) CheckProjectLimit(userID string, currentCount int) (bool, error) {
	sub, err := s.subRepo.GetSubscriptionByUserID(userID)
	if err != nil {
		return false, err
	}

	tier := sub.Tier
	if sub.Status != StatusActive {
		tier = TierFree
	}

	config, exists := Plans[tier]
	if !exists {
		return false, errors.New("unknown subscription tier")
	}

	if config.MaxProjects == -1 {
		return true, nil
	}

	return currentCount < config.MaxProjects, nil
}

func (s *Service) CheckShareLimit(userID string, currentCount int) (bool, error) {
	sub, err := s.subRepo.GetSubscriptionByUserID(userID)
	if err != nil {
		return false, err
	}

	tier := sub.Tier
	if sub.Status != StatusActive {
		tier = TierFree
	}

	config, exists := Plans[tier]
	if !exists {
		return false, errors.New("unknown subscription tier")
	}

	if config.MaxShares == -1 {
		return true, nil
	}

	return currentCount < config.MaxShares, nil
}

func (s *Service) CreateCheckoutSession(userID string, tier SubscriptionTier) (string, error) {
	var planID string
	if tier == TierSilver {
		planID = s.cfg.RazorpaySilverPlanID
	} else if tier == TierGold {
		planID = s.cfg.RazorpayGoldPlanID
	} else {
		return "", errors.New("invalid tier selected for subscription")
	}

	if planID == "" {
		return "", errors.New("razorpay plan ID not configured for this tier")
	}

	payload := RazorpaySubscriptionRequest{
		PlanID:         planID,
		TotalCount:     120, // 10 years recurring monthly
		Quantity:       1,
		CustomerNotify: 0,
		Notes: map[string]string{
			"user_id": userID,
			"tier":    string(tier),
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.razorpay.com/v1/subscriptions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(s.cfg.RazorpayKeyID, s.cfg.RazorpayKeySecret)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("razorpay API error: status %d, response: %s", resp.StatusCode, string(respBody))
	}

	var rpResp RazorpaySubscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpResp); err != nil {
		return "", err
	}

	plan, exists := Plans[tier]
	amount := 0
	if exists {
		amount = int(plan.PriceInRs * 100)
	}

	tierStr := string(tier)
	err = s.subRepo.CreateTransaction(userID, rpResp.ID, nil, amount, "created", &tierStr)
	if err != nil {
		return "", err
	}

	return rpResp.ID, nil
}

func (s *Service) VerifySubscriptionPayment(userID string, subscriptionID string, paymentID string, signature string) error {
	data := paymentID + "|" + subscriptionID
	h := hmac.New(sha256.New, []byte(s.cfg.RazorpayKeySecret))
	h.Write([]byte(data))
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != signature {
		return errors.New("invalid payment signature")
	}

	tx, err := s.subRepo.GetTransactionBySubscriptionID(subscriptionID)
	if err != nil {
		return err
	}
	if tx == nil {
		return errors.New("transaction record not found for this subscription ID")
	}

	if tx.UserID != userID {
		return errors.New("unauthorized verification request")
	}

	if tx.Status == "captured" {
		return nil
	}

	tier := SubscriptionTier(TierFree)
	if tx.BillingReason != nil {
		tier = SubscriptionTier(*tx.BillingReason)
	}

	nextMonth := time.Now().AddDate(0, 1, 3) // 1 month + 3 days grace period
	err = s.subRepo.UpdateSubscription(userID, tier, StatusActive, ProviderRazorpay, subscriptionID, &nextMonth)
	if err != nil {
		return err
	}

	capturedStatus := "captured"
	reason := "subscription.created"
	err = s.subRepo.UpdateTransaction(tx.ID, &paymentID, capturedStatus, &reason)
	if err != nil {
		fmt.Printf("failed to update transaction: %v\n", err)
	}

	return nil
}

func (s *Service) ProcessWebhook(payload []byte, headerSignature string) error {
	h := hmac.New(sha256.New, []byte(s.cfg.RazorpayWebhookSecret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	if expectedSignature != headerSignature {
		return errors.New("invalid webhook signature")
	}

	var wp RazorpayWebhookPayload
	if err := json.Unmarshal(payload, &wp); err != nil {
		return err
	}

	subID := wp.Payload.Subscription.Entity.ID
	if subID == "" && wp.Payload.Payment.Entity.SubscriptionID != "" {
		subID = wp.Payload.Payment.Entity.SubscriptionID
	}

	if subID == "" {
		return errors.New("missing subscription ID in webhook payload")
	}

	sub, err := s.subRepo.GetSubscriptionByProviderSubID(subID)
	if err != nil {
		return err
	}
	if sub == nil {
		userID := wp.Payload.Subscription.Entity.Notes["user_id"]
		if userID == "" {
			return fmt.Errorf("user not found for subscription ID: %s", subID)
		}
		sub, err = s.subRepo.GetSubscriptionByUserID(userID)
		if err != nil {
			return err
		}
	}

	switch wp.Event {
	case "subscription.charged":
		var currentPeriodEnd *time.Time
		if wp.Payload.Subscription.Entity.CurrentEnd != nil {
			t := time.Unix(*wp.Payload.Subscription.Entity.CurrentEnd, 0)
			currentPeriodEnd = &t
		} else {
			t := time.Now().AddDate(0, 1, 3)
			currentPeriodEnd = &t
		}

		tier := sub.Tier
		if wp.Payload.Subscription.Entity.Notes["tier"] != "" {
			tier = SubscriptionTier(wp.Payload.Subscription.Entity.Notes["tier"])
		}

		err = s.subRepo.UpdateSubscription(sub.UserID, tier, StatusActive, ProviderRazorpay, subID, currentPeriodEnd)
		if err != nil {
			return err
		}

		tx, err := s.subRepo.GetTransactionBySubscriptionID(subID)
		if err != nil {
			fmt.Printf("failed to lookup transaction in webhook: %v\n", err)
		}

		payID := wp.Payload.Payment.Entity.ID
		amount := wp.Payload.Payment.Entity.Amount
		reason := wp.Event

		if tx != nil && tx.Status == "created" {
			err = s.subRepo.UpdateTransaction(tx.ID, &payID, "captured", &reason)
		} else {
			err = s.subRepo.CreateTransaction(sub.UserID, subID, &payID, amount, "captured", &reason)
		}
		if err != nil {
			fmt.Printf("failed to log webhook transaction: %v\n", err)
		}

	case "subscription.cancelled":
		err = s.subRepo.UpdateSubscription(sub.UserID, sub.Tier, StatusCanceled, ProviderRazorpay, subID, sub.CurrentPeriodEnd)
		if err != nil {
			return err
		}

	case "subscription.halted":
		err = s.subRepo.UpdateSubscription(sub.UserID, sub.Tier, StatusPastDue, ProviderRazorpay, subID, sub.CurrentPeriodEnd)
		if err != nil {
			return err
		}

	case "payment.failed":
		tx, err := s.subRepo.GetTransactionBySubscriptionID(subID)
		if err != nil {
			fmt.Printf("failed to lookup transaction in webhook: %v\n", err)
		}

		payID := wp.Payload.Payment.Entity.ID
		amount := wp.Payload.Payment.Entity.Amount
		reason := wp.Event

		if tx != nil && tx.Status == "created" {
			err = s.subRepo.UpdateTransaction(tx.ID, &payID, "failed", &reason)
		} else {
			err = s.subRepo.CreateTransaction(sub.UserID, subID, &payID, amount, "failed", &reason)
		}
		if err != nil {
			fmt.Printf("failed to log webhook transaction: %v\n", err)
		}
	}

	return nil
}

func (s *Service) CancelSubscriptionCheckout(userID string, subscriptionID string) error {
	tx, err := s.subRepo.GetTransactionBySubscriptionID(subscriptionID)
	if err != nil {
		return err
	}
	if tx == nil {
		return errors.New("transaction record not found")
	}

	if tx.UserID != userID {
		return errors.New("unauthorized request")
	}

	if tx.Status == "created" {
		reason := "user_cancelled"
		err = s.subRepo.UpdateTransaction(tx.ID, nil, "failed", &reason)
		return err
	}

	return nil
}

func (s *Service) GetTransactionsByUserID(userID string) ([]SubscriptionTransaction, error) {
	return s.subRepo.GetTransactionsByUserID(userID)
}

