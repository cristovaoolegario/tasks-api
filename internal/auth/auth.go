package auth

import (
	"errors"
	"github.com/cristovaoolegario/tasks-api/internal/domain/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

type Service interface {
	GenerateJWT(user, role, id string) (string, error)
	Login(username, password string) (string, error)
	ExtractUserFromContext(ctx *gin.Context) (uint, string, error)
}

// ServiceImp used to authenticate users
type ServiceImp struct {
	secret      string
	userService service.UserService
}

func NewAuthService(secret string, userService service.UserService) *ServiceImp {
	return &ServiceImp{
		secret:      secret,
		userService: userService,
	}
}

type jwtClaims struct {
	User string `json:"username"`
	Role string `json:"role"`
	Id   string `json:"id"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token
func (s *ServiceImp) GenerateJWT(user, role, id string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &jwtClaims{
		User: user,
		Role: role,
		Id:   id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secret))

	return tokenString, err
}

// Login validates username and password and gets a valid token
func (s *ServiceImp) Login(username, password string) (string, error) {
	user, err := s.userService.FindByUsername(username)

	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if verifyPassword(user.Password, password) {
		tokenString, err := s.GenerateJWT(username, string(user.Role), strconv.Itoa(int(user.ID)))
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

func (s *ServiceImp) ExtractUserFromContext(ctx *gin.Context) (uint, string, error) {
	claimsRaw, exists := ctx.Get("user_claims")
	if !exists {
		return 0, "", errors.New("no user id claim")
	}

	claims, ok := claimsRaw.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("no user id claim")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return 0, "", errors.New("no user id claim")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", errors.New("no role claim")
	}
	userId, _ := strconv.Atoi(id)
	return uint(userId), role, nil
}
