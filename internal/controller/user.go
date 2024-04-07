package controller

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/dto"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/domain/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	service service.UserService
}

func NewUserController(service service.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Add a new user to the system
// @Tags user
// @Accept  json
// @Produce  json
// @Param   user  body     dto.User  true  "User to create"
// @Success 201  {object}  dto.User
// @Failure 400  {object}  map[string]interface{}  "Input validation error"
// @Failure 500  {object}  map[string]interface{}  "Internal server error"
// @Router /api/users [post]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (uc *UserController) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.service.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &dto.User{
		ID:       user.ID,
		Username: user.Username,
		Role:     string(user.Role),
	})
}

// GetUser godoc
// @Summary Get a user by username
// @Description Retrieve user details by username
// @Tags user
// @Accept  json
// @Produce  json
// @Param   username  path      string  true  "Username"
// @Success 200  {object}  dto.User
// @Failure 404  {object}  map[string]interface{}  "User not found"
// @Router /api/users/{username} [get]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (uc *UserController) GetUser(c *gin.Context) {
	userName := c.Param("username")
	user, err := uc.service.FindByUsername(userName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, &dto.User{
		ID:       user.ID,
		Username: user.Username,
		Role:     string(user.Role),
	})
}
