package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenAuthMiddleware(t *testing.T) {
	secret := "secret"
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TokenAuthMiddleware(secret))

	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Authenticated"})
	})

	t.Run("Should return status OK When token is valid", func(t *testing.T) {
		validToken, _ := NewAuthService(secret, nil).GenerateJWT("new user", "manager", "1")

		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Should return statys Unauthorized When token is invalid ", func(t *testing.T) {
		invalidToken := "invalid_token"
		req := httptest.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+invalidToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Should return status When no token is provided", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRoleManagerMiddleware(t *testing.T) {
	secret := "secret"
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TokenAuthMiddleware(secret))

	router.GET("/manager-protected", RoleManagerMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Authenticated as manager"})
	})

	t.Run("Should return status OK When the role is manager", func(t *testing.T) {
		validToken, _ := NewAuthService(secret, nil).GenerateJWT("new user", "manager", "1")
		req := httptest.NewRequest("GET", "/manager-protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Should return status StatusForbidden When the role is technician", func(t *testing.T) {
		validToken, _ := NewAuthService(secret, nil).GenerateJWT("new user", "technician", "1")
		req := httptest.NewRequest("GET", "/manager-protected", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Should return status StatusUnauthorized When there is no claims", func(t *testing.T) {
		w := performRequestWithoutClaims(router, "GET", "/manager-protected")
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func performRequestWithoutClaims(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req

	r.ServeHTTP(w, req)
	return w
}
