package controller

import (
	"bytes"
	"encoding/json"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/domain/service"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserController(t *testing.T) {

	repository := repository.NewMockUserRepository()
	service := service.NewUserService(repository)
	userController := NewUserController(service)

	t.Run("Should return 201 When able to create user", func(t *testing.T) {
		w := performRequest(userController.CreateUser, http.MethodPost, "/users", "", model.User{Role: model.Manager, Username: "new_user", Password: "test_password"})

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Should return 200 and user When user exists", func(t *testing.T) {
		w := performRequest(userController.GetUser, http.MethodGet, "/users", "new_user", nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Should return 404 When user doesn't exists", func(t *testing.T) {
		w := performRequest(userController.GetUser, http.MethodGet, "/users", "no_existing_user", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

func performRequest(handlerFunc gin.HandlerFunc, method, path, username string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			panic(err) // Handle error accordingly
		}
	}

	req, _ := http.NewRequest(method, path, bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	if username != "" {
		c.AddParam("username", username)
	}
	c.Request = req

	handlerFunc(c)
	return w
}
