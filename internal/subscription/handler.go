package subscription

import (
	"github.com/gin-gonic/gin"
	"github.com/shubham071122/collab/internal/response"
)

type Handler struct {
	subService *Service
}

func NewHandler(subService *Service) *Handler {
	return &Handler{
		subService: subService,
	}
}

func (h *Handler) GetPlans(c *gin.Context) {
	orderedTiers := []SubscriptionTier{TierFree, TierSilver, TierGold}
	planList := make([]PlanConfig, 0, len(orderedTiers))

	for _, tier := range orderedTiers {
		if plan, exists := Plans[tier]; exists {
			planList = append(planList, plan)
		}
	}

	response.JSON(c, response.StatusOK, "Plans retrieved successfully", planList, nil)
}

func (h *Handler) GetSubscription(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.JSON(c, response.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context")
		return
	}

	sub, err := h.subService.GetSubscriptionByUserID(userID)

	if err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to retrieve subscription", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Subscription retrieved successfully", sub, nil)
}

func (h *Handler) CreateCheckout(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.JSON(c, response.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context")
		return
	}

	var req struct {
		Tier SubscriptionTier `json:"tier"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "Invalid request body", nil, err.Error())
		return
	}

	subID, err := h.subService.CreateCheckoutSession(userID, req.Tier)
	if err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to create subscription checkout session", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Checkout session created successfully", gin.H{
		"subscription_id": subID,
	}, nil)
}

func (h *Handler) VerifyPayment(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.JSON(c, response.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context")
		return
	}

	var req struct {
		SubscriptionID    string `json:"subscription_id"`
		RazorpayPaymentID string `json:"razorpay_payment_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "Invalid request body", nil, err.Error())
		return
	}

	err := h.subService.VerifySubscriptionPayment(userID, req.SubscriptionID, req.RazorpayPaymentID, req.RazorpaySignature)
	if err != nil {
		response.JSON(c, response.StatusBadRequest, "Payment verification failed", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Payment verified and subscription activated successfully", nil, nil)
}

func (h *Handler) CancelCheckout(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.JSON(c, response.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context")
		return
	}

	var req struct {
		SubscriptionID string `json:"subscription_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, response.StatusBadRequest, "Invalid request body", nil, err.Error())
		return
	}

	err := h.subService.CancelSubscriptionCheckout(userID, req.SubscriptionID)
	if err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to cancel checkout session", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Checkout session cancelled successfully", nil, nil)
}

func (h *Handler) GetTransactions(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		response.JSON(c, response.StatusUnauthorized, "Unauthorized", nil, "User ID not found in context")
		return
	}

	txs, err := h.subService.GetTransactionsByUserID(userID)
	if err != nil {
		response.JSON(c, response.StatusInternalServerError, "Failed to retrieve transactions", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Transactions retrieved successfully", txs, nil)
}



func (h *Handler) Webhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		response.JSON(c, response.StatusBadRequest, "Failed to read request body", nil, err.Error())
		return
	}

	signature := c.GetHeader("X-Razorpay-Signature")
	if signature == "" {
		response.JSON(c, response.StatusBadRequest, "Missing X-Razorpay-Signature header", nil, nil)
		return
	}

	err = h.subService.ProcessWebhook(payload, signature)
	if err != nil {
		response.JSON(c, response.StatusBadRequest, "Webhook processing failed", nil, err.Error())
		return
	}

	response.JSON(c, response.StatusOK, "Webhook processed successfully", nil, nil)
}

