package auth

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/domain/service"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestService_GenerateJWT(t *testing.T) {
	secret := "testing-secret"
	service := NewAuthService(secret, service.NewUserService(&repository.MockUserRepository{}))

	username := "testUser"
	role := "technician"
	id := "1"

	tokenString, err := service.GenerateJWT(username, role, id)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Fatalf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(*jwtClaims); ok && token.Valid {
		assert.Equal(t, username, claims.User)
		assert.Equal(t, role, claims.Role)
		assert.Equal(t, id, claims.Id)
		assert.LessOrEqual(t, time.Now().Unix(), claims.ExpiresAt)
	}
}

func TestService_Login(t *testing.T) {
	mockRepo := repository.NewMockUserRepository()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	mockUser := &model.User{Username: "testuser", Password: string(hashedPassword)}
	mockRepo.Users[1] = mockUser

	service := &ServiceImp{
		userService: service.NewUserService(mockRepo),
	}

	t.Run("Should return a token When login is valid", func(t *testing.T) {
		token, err := service.Login("testuser", "password") // Assuming "password" is the correct password
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Should return error When login is invalid", func(t *testing.T) {
		emptyToken, err := service.Login("testuser", "wrongpassword")
		assert.Error(t, err)
		assert.Empty(t, emptyToken)
	})

	t.Run("Should return error When user doesn't exists", func(t *testing.T) {
		emptyToken, err := service.Login("nonexistentuser", "password")
		assert.Error(t, err)
		assert.Empty(t, emptyToken)
	})
}

func TestServiceImp_ExtractUserIdFromContext(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("Should extract userid When claims are set", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		expectedUserID := uint(123)
		claims := jwt.MapClaims{
			"id": strconv.Itoa(int(expectedUserID)),
		}
		c.Set("user_claims", claims)

		secret := "testing-secret"
		service := NewAuthService(secret, service.NewUserService(&repository.MockUserRepository{}))

		userID, err := service.ExtractUserIdFromContext(c)

		assert.NoError(t, err)
		assert.Equal(t, expectedUserID, userID)
	})

	t.Run("Should return error When claims weren't set", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		secret := "testing-secret"
		service := NewAuthService(secret, service.NewUserService(&repository.MockUserRepository{}))

		userID, err := service.ExtractUserIdFromContext(c)

		assert.Error(t, err)
		assert.Equal(t, uint(0), userID)
	})

}
