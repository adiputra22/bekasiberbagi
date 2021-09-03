package handler

import (
	"bekasiberbagi/auth"
	"bekasiberbagi/response"
	"bekasiberbagi/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := response.APIResponseValidationFailed("Register failed", http.StatusUnprocessableEntity, err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := response.APIResponseFailed("Register failed", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := response.APIResponseFailed(err.Error(), http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userFormatter := user.FormatUser(newUser, token)

	response := response.APIResponseSuccess("User has been registered", http.StatusOK, userFormatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := response.APIResponseValidationFailed("Login failed", http.StatusUnprocessableEntity, err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedUser, err := h.userService.Login(input)

	if err != nil {
		response := response.APIResponseFailed(err.Error(), http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedUser.ID)
	if err != nil {
		response := response.APIResponseFailed(err.Error(), http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedUser, token)

	response := response.APIResponseSuccess("Success login", http.StatusOK, formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) IsEmailAvailability(c *gin.Context) {
	var input user.CheckEmailAvailabilityInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := response.APIResponseValidationFailed("Input failed", http.StatusUnprocessableEntity, err)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedUser, err := h.userService.IsEmailAvailability(input)

	if err != nil {
		response := response.APIResponseFailed(err.Error(), http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := response.APIResponseSuccess("Success. User not found", http.StatusOK, loggedUser)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := response.APIResponseFailedWithData(err.Error(), http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)

	userId := currentUser.ID

	path := fmt.Sprintf("uploads/images/%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := response.APIResponseFailedWithData(err.Error(), http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userId, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := response.APIResponseFailedWithData(err.Error(), http.StatusBadRequest, data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := response.APIResponseSuccess("Success upload avatar", http.StatusOK, data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	formatter := user.FormatUser(currentUser, "")

	response := response.APIResponseSuccess("Success upload avatar", http.StatusOK, formatter)

	c.JSON(http.StatusOK, response)
}
