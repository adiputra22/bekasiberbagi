package campaign

type Service interface {
	FindCampaigns(userId int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) FindCampaigns(userId int) ([]Campaign, error) {
	campaigns, err := s.repository.FindAll()

	if userId != 0 {
		campaigns, err = s.repository.FindByUserID(userId)
	}

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
