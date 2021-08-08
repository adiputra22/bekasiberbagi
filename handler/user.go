package handler

import (
	"bekasiberbagi/response"
	"bekasiberbagi/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
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

	userFormatter := user.FormatUser(newUser, "test")

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

	formatter := user.FormatUser(loggedUser, "token")

	response := response.APIResponseSuccess("Success login", http.StatusOK, formatter)

	c.JSON(http.StatusOK, response)
}
