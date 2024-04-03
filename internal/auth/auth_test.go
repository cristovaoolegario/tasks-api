package auth

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestService_GenerateJWT(t *testing.T) {
	secret := "testing-secret"
	service := NewAuthService(secret, &repository.MockUserRepository{})

	username := "testUser"
	role := "technician"

	tokenString, err := service.GenerateJWT(username, role)
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
		if claims.User != username {
			t.Errorf("Expected username %v, got %v", username, claims.User)
		}
		if claims.Role != role {
			t.Errorf("Expected role %v, got %v", role, claims.Role)
		}

		now := time.Now().Unix()
		if claims.ExpiresAt < now {
			t.Errorf("Token should not be expired")
		}

	} else {
		t.Fatalf("Failed to parse token: %v", err)
	}
}

func TestService_Login(t *testing.T) {
	// Setup
	mockRepo := repository.NewMockUserRepository()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	mockUser := &model.User{Username: "testuser", Password: string(hashedPassword)}
	mockRepo.Users[1] = mockUser

	service := &Service{
		repo: mockRepo,
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
