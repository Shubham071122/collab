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
