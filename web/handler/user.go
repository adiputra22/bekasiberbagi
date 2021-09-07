package handler

import (
	"bekasiberbagi/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUsers()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"users": users})
}

func (h *userHandler) Create(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}

func (h *userHandler) Store(c *gin.Context) {
	var form user.FormCreateInput

	err := c.ShouldBind(&form)
	if err != nil {
		form.Error = err
		c.HTML(http.StatusInternalServerError, "create.html", form)
		return
	}

	_, err = h.userService.StoreFromForm(form)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/web/users")
}

func (h *userHandler) Edit(c *gin.Context) {
	idParam, _ := strconv.Atoi(c.Param("id"))

	registeredUser, err := h.userService.GetUserById(idParam)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := user.FormUpdateInput{}
	input.ID = registeredUser.ID
	input.Name = registeredUser.Name
	input.Email = registeredUser.Email
	input.Occupation = registeredUser.Occupation
	input.Error = nil

	c.HTML(http.StatusOK, "edit.html", input)
}

func (h *userHandler) Update(c *gin.Context) {
	idParam, _ := strconv.Atoi(c.Param("id"))

	userExists, err := h.userService.GetUserById(idParam)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	var form user.FormUpdateInput

	err = c.ShouldBind(&form)
	if err != nil {
		form.Error = err
		c.HTML(http.StatusInternalServerError, "edit.html", userExists)
		return
	}

	form.ID = idParam

	_, err = h.userService.UpdateFromForm(form)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/web/users")
}

func (h *userHandler) EditAvatar(c *gin.Context) {
	idParam, _ := strconv.Atoi(c.Param("id"))

	registeredUser, err := h.userService.GetUserById(idParam)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := user.FormUpdateAvatar{}
	input.ID = registeredUser.ID
	input.Name = registeredUser.Name
	input.AvatarFileName = registeredUser.AvatarFileName
	input.Error = nil

	c.HTML(http.StatusOK, "edit_avatar.html", input)
}

func (h *userHandler) UpdateAvatar(c *gin.Context) {
	idParam, _ := strconv.Atoi(c.Param("id"))

	userExists, err := h.userService.GetUserById(idParam)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	file, err := c.FormFile("avatar")

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	path := fmt.Sprintf("uploads/images/%d-%s", idParam, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	var form user.FormUpdateAvatar

	form.ID = idParam
	form.AvatarFileName = path
	form.Name = userExists.Name

	_, err = h.userService.UpdateAvatarFromForm(form)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/web/users")
}
