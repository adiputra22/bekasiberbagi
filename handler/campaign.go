package handler

import (
	"bekasiberbagi/campaign"
	"bekasiberbagi/response"
	"bekasiberbagi/user"
	"fmt"
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

func (h *CampaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := response.APIResponseValidationFailed("Create campaign failed coz input", http.StatusUnprocessableEntity, err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	input.User = c.MustGet("currentUser").(user.User)

	campaignCreated, err := h.service.CreateCampaign(input)

	if err != nil {
		response := response.APIResponseValidationFailed("Create campaign failed coz service", http.StatusUnprocessableEntity, err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := response.APIResponseSuccess("Campaign detail", http.StatusOK, campaign.FormatCampaign(campaignCreated))
	c.JSON(http.StatusOK, response)
}

func (h *CampaignHandler) UpdateCampaign(c *gin.Context) {
	var inputUri campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&inputUri)
	if err != nil {
		response := response.APIResponseFailed("Update campaign failed coz error uri", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var input campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&input)

	if err != nil {
		response := response.APIResponseValidationFailed("Update campaign failed coz input", http.StatusUnprocessableEntity, err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	input.User = c.MustGet("currentUser").(user.User)

	campaignCreated, err := h.service.UpdateCampaign(inputUri, input)

	if err != nil {
		data := gin.H{"error": err.Error()}
		response := response.APIResponseFailedWithData("Update campaign failed coz service", http.StatusUnprocessableEntity, data)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := response.APIResponseSuccess("Update campaign success", http.StatusOK, campaign.FormatCampaign(campaignCreated))
	c.JSON(http.StatusOK, response)
}

func (h *CampaignHandler) CreateCampaignImage(c *gin.Context) {
	file, err := c.FormFile("campaign_image")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := response.APIResponseFailedWithData(err.Error(), http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	userId := currentUser.ID

	path := fmt.Sprintf("uploads/campaign/%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := response.APIResponseFailedWithData(err.Error(), http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var createCampaignImageInput campaign.CreateCampaignImageInput

	err = c.ShouldBind(&createCampaignImageInput)

	createCampaignImageInput.User = currentUser

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := response.APIResponseFailedWithData(err.Error(), http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.CreateCampaignImage(createCampaignImageInput, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := response.APIResponseFailedWithData(err.Error(), http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := response.APIResponseSuccess("Success upload campaign image", http.StatusOK, data)
	c.JSON(http.StatusOK, response)
}
