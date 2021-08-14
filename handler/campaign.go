package handler

import (
	"bekasiberbagi/campaign"
	"bekasiberbagi/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *CampaignHandler {
	return &CampaignHandler{service}
}

func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userId)

	if err != nil {
		response := response.APIResponseFailed("Get Campaigns failed", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := response.APIResponseSuccess("List of campaigns", http.StatusOK, campaign.FormatCampaigns(campaigns))

	c.JSON(http.StatusOK, response)
}
