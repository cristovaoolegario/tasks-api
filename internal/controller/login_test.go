package controller_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/cristovaoolegario/tasks-api/internal/controller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	t.Run("Should return 200 and token", func(t *testing.T) {
		loginController := controller.NewLoginController(
			&auth.MockAuthService{
				LoginMock: func(username, password string) (string, error) {
					return "mocked_token", nil
				},
			})

		r := gin.Default()

		r.POST("/login", loginController.LoginHandler)
		req, _ := http.NewRequest("POST", "/login", nil)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponseBody := `{"token":"mocked_token"}`
		assert.JSONEq(t, expectedResponseBody, rec.Body.String())
	})

	t.Run("Should return 401", func(t *testing.T) {
		loginController := controller.NewLoginController(
			&auth.MockAuthService{
				LoginMock: func(username, password string) (string, error) {
					return "", errors.New("some error")
				},
			})

		r := gin.Default()

		r.POST("/login", loginController.LoginHandler)
		req, _ := http.NewRequest("POST", "/login", nil)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
