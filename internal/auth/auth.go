package auth

import (
	"errors"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Service used to authenticate users
type Service struct {
	secret string
	repo   repository.UserRepository
}

func NewAuthService(secret string, userRepository repository.UserRepository) *Service {
	return &Service{
		secret: secret,
		repo:   userRepository,
	}
}

type jwtClaims struct {
	User string `json:"username"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token
func (s *Service) GenerateJWT(user, role string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwtClaims{
		User: user,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))

	return tokenString, err
}

// Login validates username and password and gets a valid token
func (s *Service) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)

	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if verifyPassword(user.Password, password) {
		tokenString, err := s.GenerateJWT(username, string(user.Role))
		if err != nil {
			return "", errors.New("error generating token")
		}
		return tokenString, nil
	}
	return "", errors.New("invalid username or password")
}

func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
