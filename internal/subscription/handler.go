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
