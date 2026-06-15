package subscription

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserSubscription struct {
	ID                     string             `json:"id"`
	UserID                 string             `json:"user_id"`
	Tier                   SubscriptionTier   `json:"tier"`
	Status                 SubscriptionStatus `json:"status"`
	Provider               PaymentProvider    `json:"provider"`
	ProviderSubscriptionID *string            `json:"provider_subscription_id"`
	CurrentPeriodEnd       *time.Time         `json:"current_period_end"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateDefaultSubscription(userID string) error {
	query := `
		INSERT INTO user_subscriptions (
			id,
			user_id,
			tier,
			status,
			provider
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(
		query,
		uuid.New().String(),
		userID,
		TierFree,
		StatusActive,
		ProviderNone,
	)
	return err
}

func (r *Repository) GetSubscriptionByUserID(userID string) (*UserSubscription, error) {
	query := `
		SELECT id, user_id, tier, status, provider, provider_subscription_id, current_period_end, created_at, updated_at
		FROM user_subscriptions
		WHERE user_id = $1
	`
	var sub UserSubscription
	err := r.db.QueryRow(query, userID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.Tier,
		&sub.Status,
		&sub.Provider,
		&sub.ProviderSubscriptionID,
		&sub.CurrentPeriodEnd,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return &UserSubscription{
			Tier:   TierFree,
			Status: StatusActive,
		}, nil
	} else if err != nil {
		return nil, err
	}
	return &sub, nil
}

type SubscriptionTransaction struct {
	ID             string     `json:"id"`
	UserID                 string     `json:"user_id"`
	SubscriptionID string     `json:"subscription_id"`
	PaymentID      *string    `json:"payment_id"`
	Amount         int        `json:"amount"`
	Status         string     `json:"status"`
	BillingReason  *string    `json:"billing_reason"`
	CreatedAt      time.Time  `json:"created_at"`
}

func (r *Repository) GetSubscriptionByProviderSubID(providerSubID string) (*UserSubscription, error) {
	query := `
		SELECT id, user_id, tier, status, provider, provider_subscription_id, current_period_end, created_at, updated_at
		FROM user_subscriptions
		WHERE provider_subscription_id = $1
	`
	var sub UserSubscription
	err := r.db.QueryRow(query, providerSubID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.Tier,
		&sub.Status,
		&sub.Provider,
		&sub.ProviderSubscriptionID,
		&sub.CurrentPeriodEnd,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *Repository) UpdateSubscription(userID string, tier SubscriptionTier, status SubscriptionStatus, provider PaymentProvider, providerSubID string, currentPeriodEnd *time.Time) error {
	query := `
		UPDATE user_subscriptions
		SET tier = $1, status = $2, provider = $3, provider_subscription_id = $4, current_period_end = $5, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $6
	`
	_, err := r.db.Exec(query, tier, status, provider, providerSubID, currentPeriodEnd, userID)
	return err
}

func (r *Repository) CreateTransaction(userID string, subscriptionID string, paymentID *string, amount int, status string, billingReason *string) error {
	query := `
		INSERT INTO subscription_transactions (
			id,
			user_id,
			subscription_id,
			payment_id,
			amount,
			status,
			billing_reason
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(
		query,
		uuid.New().String(),
		userID,
		subscriptionID,
		paymentID,
		amount,
		status,
		billingReason,
	)
	return err
}

func (r *Repository) GetTransactionBySubscriptionID(subID string) (*SubscriptionTransaction, error) {
	query := `
		SELECT id, user_id, subscription_id, payment_id, amount, status, billing_reason, created_at
		FROM subscription_transactions
		WHERE subscription_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var tx SubscriptionTransaction
	err := r.db.QueryRow(query, subID).Scan(
		&tx.ID,
		&tx.UserID,
		&tx.SubscriptionID,
		&tx.PaymentID,
		&tx.Amount,
		&tx.Status,
		&tx.BillingReason,
		&tx.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *Repository) UpdateTransaction(id string, paymentID *string, status string, billingReason *string) error {
	query := `
		UPDATE subscription_transactions
		SET payment_id = $1, status = $2, billing_reason = $3
		WHERE id = $4
	`
	_, err := r.db.Exec(query, paymentID, status, billingReason, id)
	return err
}

func (r *Repository) GetTransactionsByUserID(userID string) ([]SubscriptionTransaction, error) {
	query := `
		SELECT id, user_id, subscription_id, payment_id, amount, status, billing_reason, created_at
		FROM subscription_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []SubscriptionTransaction
	for rows.Next() {
		var tx SubscriptionTransaction
		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.SubscriptionID,
			&tx.PaymentID,
			&tx.Amount,
			&tx.Status,
			&tx.BillingReason,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}



