package subscription

type PlanConfig struct {
	ID          SubscriptionTier `json:"id"`
	Name        string           `json:"name"`
	PriceInRs   float64          `json:"price_in_rs"`
	MaxProjects int              `json:"max_projects"`
	MaxShares   int              `json:"max_shares"`
	Description string           `json:"description"`
	Features    []string         `json:"features"`
}

var Plans = map[SubscriptionTier]PlanConfig{
	TierFree: {
		ID:          TierFree,
		Name:        "Free",
		PriceInRs:   0.0,
		MaxProjects: 2,
		MaxShares:   2,
		Description: "Get started with up to 2 projects and share with 2 collaborators.",
		Features: []string{
			"Up to 2 projects",
			"Up to 2 collaborators per project",
			"Infinite collaborative canvas",
			"Real-time socket sync",
		},
	},
	TierSilver: {
		ID:          TierSilver,
		Name:        "Silver",
		PriceInRs:   5.0,
		MaxProjects: 5,
		MaxShares:   5,
		Description: "Scale up with up to 5 projects and share with 5 collaborators.",
		Features: []string{
			"Up to 5 projects",
			"Up to 5 collaborators per project",
			"Infinite collaborative canvas",
			"Real-time socket sync",
			"Priority project access",
		},
	},
	TierGold: {
		ID:          TierGold,
		Name:        "Gold",
		PriceInRs:   10.0,
		MaxProjects: -1, // Unlimited
		MaxShares:   -1, // Unlimited
		Description: "Unlock full capability with unlimited projects and collaborators.",
		Features: []string{
			"Unlimited projects",
			"Unlimited collaborators",
			"Infinite collaborative canvas",
			"Real-time socket sync",
			"Priority premium support",
		},
	},
}
