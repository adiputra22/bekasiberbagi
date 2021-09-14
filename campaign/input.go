package campaign

import "bekasiberbagi/user"

type GetCampaignDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	Perks            string `json:"perks"`
	User             user.User
}

type CreateCampaignImageInput struct {
	CampaignId int   `form:"campaign_id" binding:"required"`
	IsPrimary  *bool `form:"is_primary" binding:"required"`
	User       user.User
}

type FormCreateCampaignInput struct {
	ID               int
	Name             string `form:"name" binding:"required"`
	ShortDescription string `form:"short_description" binding:"required"`
	Description      string `form:"description" binding:"required"`
	GoalAmount       int    `form:"goal_amount" binding:"required"`
	Perks            string `form:"perks" binding:"required"`
	UserID           int    `form:"user_id" binding:"required"`
	Error            error
	Users            []user.User
}

type FormUpdateImage struct {
	ID    int
	Name  string `form:"name" binding:"required"`
	Image string `file:"campaign_image"`
	Error error
}
