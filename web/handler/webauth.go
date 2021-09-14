package handler

import (
	"bekasiberbagi/user"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type webAuthHandler struct {
	userService user.Service
}

func NewWebAuthHandler(userService user.Service) *webAuthHandler {
	return &webAuthHandler{
		userService: userService,
	}
}

func (h *webAuthHandler) LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *webAuthHandler) LoginAction(c *gin.Context) {
	var input user.WebLoginInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	user, err := h.userService.WebLogin(input)
	if err != nil || user.Role != "admin" {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	session := sessions.Default(c)
	session.Set("userId", user.ID)
	session.Set("userName", user.Name)
	session.Save()

	c.Redirect(http.StatusFound, "/web/users")
}
