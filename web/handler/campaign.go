package handler

import (
	"bekasiberbagi/campaign"
	"bekasiberbagi/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	campaignService campaign.Service
	userService     user.Service
}

func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignHandler {
	return &campaignHandler{
		campaignService: campaignService,
		userService:     userService,
	}
}

func (h *campaignHandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetAllCampaigns()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

func (h *campaignHandler) Create(c *gin.Context) {
	users, err := h.userService.GetAllUsers()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	form := campaign.FormCreateCampaignInput{}
	form.Users = users

	c.HTML(http.StatusOK, "campaign_create.html", form)
}

func (h *campaignHandler) Store(c *gin.Context) {
	var form campaign.FormCreateCampaignInput

	err := c.ShouldBind(&form)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	_, err = h.campaignService.CreateFromForm(form)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/web/campaigns")
}

func (h *campaignHandler) Edit(c *gin.Context) {
	idParam, _ := strconv.Atoi(c.Param("id"))

	campaignRegistered, err := h.campaignService.GetCampaignByIntId(idParam)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	var input campaign.FormCreateCampaignInput
	input.ID = campaignRegistered.ID
	input.Name = campaignRegistered.Name
	input.ShortDescription = campaignRegistered.ShortDescription
	input.Description = campaignRegistered.Description
	input.GoalAmount = campaignRegistered.GoalAmount
	input.Perks = campaignRegistered.Perks
	input.UserID = campaignRegistered.UserID
	input.Error = nil

	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	input.Users = users

	c.HTML(http.StatusOK, "campaign_edit.html", input)
}

func (h *campaignHandler) Update(c *gin.Context) {
	var form campaign.FormCreateCampaignInput

	err := c.ShouldBind(&form)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	idParam, _ := strconv.Atoi(c.Param("id"))
	form.ID = idParam

	_, err = h.campaignService.UpdateFromForm(form)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/web/campaigns")
}

func (h *campaignHandler) FormUploadImage(c *gin.Context) {
	idParam, _ := strconv.Atoi(c.Param("id"))

	campaignRegistered, err := h.campaignService.GetCampaignByIntId(idParam)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	var input campaign.FormUpdateImage
	input.ID = campaignRegistered.ID
	input.Name = campaignRegistered.Name

	primaryImage := ""

	if len(campaignRegistered.CampaignImages) > 0 {
		primaryImage = campaignRegistered.CampaignImages[0].FileName
	}

	input.Image = primaryImage
	input.Error = nil

	c.HTML(http.StatusOK, "campaign_image.html", input)
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var form campaign.FormUpdateImage

	err := c.ShouldBind(&form)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	file, err := c.FormFile("campaign_image")

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	idParam, _ := strconv.Atoi(c.Param("id"))

	campaignRegistered, err := h.campaignService.GetCampaignByIntId(idParam)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	form.ID = idParam

	userId := campaignRegistered.UserID

	path := fmt.Sprintf("uploads/campaign/%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	_, err = h.campaignService.UploadImageFromForm(form, path)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/web/campaigns")
}
