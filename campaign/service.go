package campaign

import (
	"errors"
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userId int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputUri GetCampaignDetailInput, input CreateCampaignInput) (Campaign, error)

	CreateCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) GetCampaigns(userId int) ([]Campaign, error) {
	campaigns, err := s.repository.FindAll()

	if userId != 0 {
		campaigns, err = s.repository.FindByUserID(userId)
	}

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}

func (s *service) GetCampaignById(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindById(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.Perks = input.Perks
	campaign.UserID = input.User.ID
	campaign.BackerCount = 0
	campaign.CreatedAt = time.Now()
	campaign.UpdatedAt = time.Now()

	userSlug := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(userSlug)

	campaign, err := s.repository.Save(campaign)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) UpdateCampaign(inputUri GetCampaignDetailInput, input CreateCampaignInput) (Campaign, error) {
	singleCampaign, err := s.repository.FindById(inputUri.ID)

	if err != nil {
		return singleCampaign, err
	}

	if singleCampaign.UserID != input.User.ID {
		return singleCampaign, errors.New("USER UNAUTHORIZED TO EDIT THIS CAMPAIGN")
	}

	updatedCampaign := Campaign{}
	updatedCampaign.ID = inputUri.ID
	updatedCampaign.Name = input.Name
	updatedCampaign.ShortDescription = input.ShortDescription
	updatedCampaign.Description = input.Description
	updatedCampaign.GoalAmount = input.GoalAmount
	updatedCampaign.Perks = input.Perks
	updatedCampaign.UserID = input.User.ID
	updatedCampaign.BackerCount = 0
	updatedCampaign.CreatedAt = singleCampaign.CreatedAt
	updatedCampaign.UpdatedAt = time.Now()
	updatedCampaign.Slug = singleCampaign.Slug

	resultCampaign, err := s.repository.Update(updatedCampaign)

	if err != nil {
		return resultCampaign, err
	}

	return resultCampaign, nil
}

func (s *service) CreateCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	campaign, err := s.repository.FindById(input.CampaignId)

	if err != nil {
		return CampaignImage{}, err
	}

	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("USER UNAUTHORIZED TO UPLOAD IMAGE THIS CAMPAIGN")
	}

	isPrimary := 0

	if *input.IsPrimary {
		isPrimary = 1

		_, err := s.repository.MarkImageToNonPrimary(input.CampaignId)

		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignId
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation
	campaignImage.CreatedAt = time.Now()
	campaignImage.UpdatedAt = time.Now()

	resultCampaignImage, err := s.repository.CreateImage(campaignImage)

	if err != nil {
		return resultCampaignImage, err
	}

	return resultCampaignImage, nil
}