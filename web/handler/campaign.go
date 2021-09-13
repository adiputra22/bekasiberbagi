package handler

import (
	"bekasiberbagi/campaign"
	"bekasiberbagi/user"
	"net/http"

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
