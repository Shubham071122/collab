package subscription

import "errors"

var ErrLimitExceeded = errors.New("LIMIT_EXCEEDED")

type SubscriptionTier string

const (
	TierFree   SubscriptionTier = "free"
	TierSilver SubscriptionTier = "silver"
	TierGold   SubscriptionTier = "gold"
)

type SubscriptionStatus string

const (
	StatusActive     SubscriptionStatus = "active"
	StatusPastDue    SubscriptionStatus = "past_due"
	StatusCanceled   SubscriptionStatus = "canceled"
	StatusIncomplete SubscriptionStatus = "incomplete"
)

type PaymentProvider string

const (
	ProviderNone     PaymentProvider = "none"
	ProviderRazorpay PaymentProvider = "razorpay"
	ProviderStripe   PaymentProvider = "stripe"
)
