package handler

import (
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
