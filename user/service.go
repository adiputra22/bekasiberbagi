package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailability(input CheckEmailAvailabilityInput) (bool, error)
	SaveAvatar(Id int, fileLocation string) (User, error)
	GetUserById(Id int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository: repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("User not found")
	}

	errPassword := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if errPassword != nil {
		return user, errPassword
	}

	return user, nil
}

func (s *service) IsEmailAvailability(input CheckEmailAvailabilityInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return false, err
	}

	if user.ID > 0 {
		return false, errors.New("User found")
	}

	return true, nil
}

func (s *service) SaveAvatar(Id int, fileLocation string) (User, error) {
	user, err := s.repository.FindById(Id)

	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation

	updatedUser, err := s.repository.Update(user)

	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetUserById(Id int) (User, error) {
	user, err := s.repository.FindById(Id)

	if err != nil {
		return user, err
	}

	return user, nil
}
