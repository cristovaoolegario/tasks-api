package controller

import (
	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginController struct {
	authService auth.Service
}

func NewLoginController(authService auth.Service) *LoginController {
	return &LoginController{
		authService: authService,
	}
}

// LoginHandler
// @Summary User login
// @Description Perform a user login
// @Tags authentication
// @Accept  mpfd
// @Produce  json
// @Param   username     formData    string true "Username"
// @Param   password     formData    string true "Password"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func (s *LoginController) LoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	token, err := s.authService.Login(username, password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
