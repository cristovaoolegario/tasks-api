package controller

import (
	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
	"net/http"
)

type LoginController struct {
	userRepo    repository.UserRepository
	authService *auth.Service
}

func NewLoginController(
	userRepo repository.UserRepository,
	authService *auth.Service) *LoginController {
	return &LoginController{
		userRepo:    userRepo,
		authService: authService,
	}
}

func (s *LoginController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	token, err := s.authService.Login(username, password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Write([]byte(token))
	return
}
