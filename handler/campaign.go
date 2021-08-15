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

func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := response.APIResponseFailed("Error uri", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignById(input)
	if err != nil {
		response := response.APIResponseFailed("Error when get detail", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := response.APIResponseSuccess("Campaign detail", http.StatusOK, campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}
