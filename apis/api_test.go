package apis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"main/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест для обновления токена
func TestRefreshToken(t *testing.T) {
	router := gin.Default()
	router.GET("/refresh_token", RefreshToken)
	// Создаем mock запрос с токеном(изза ограничений на использование токенов нужно менять каждый тест)
	req, _ := http.NewRequest("GET", "/refresh_token?token=valid_refresh_token", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	// Проверка, что запрос прошел успешно
	assert.Equal(t, http.StatusOK, resp.Code)
	// Проверка содержимого
	var response models.JsonResponce
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.NotEmpty(t, response.Access)
	assert.NotEmpty(t, response.Refresh)
	assert.NotEmpty(t, response.Uuid)
}

// Тест для получения токена по uuid
func TestGetToken(t *testing.T) {
	router := gin.Default()
	router.GET("/getToken", GetToken)
	// Создаем mock запрос с uuid(изза ограничений на использование токенов нужно менять каждый тест)
	req, _ := http.NewRequest("GET", "/getToken?uuid=99a79549-b0bc-49f4-8f36-2c5be1e5e805", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	// Проверка, что запрос прошел успешно
	assert.Equal(t, http.StatusOK, resp.Code)
	// Проверка содержимого
	var response models.JsonResponce
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.NotEmpty(t, response.Access)
	assert.NotEmpty(t, response.Refresh)
	assert.Equal(t, "uUiD", response.Uuid)
}

// Тест для корректной регистрации
func TestRegisterUser(t *testing.T) {
	router := gin.Default()
	router.GET("/register", RegisterUser)
	// Создаем mock запрос с данными пользователя(изза ограничений на уникальность нужно менять каждый тест)
	req, _ := http.NewRequest("GET", "/register?username=fuhseruseffrniy1223&email=fmail23@mafil.com", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	// Проверка, что запрос прошел успешно
	assert.Equal(t, http.StatusOK, resp.Code)
	// Проверка содержимого
	assert.NotEmpty(t, resp.Body.String())
}

// Тестирование ошибок при регистрации
func TestRegisterUserWithError(t *testing.T) {
	router := gin.Default()
	router.GET("/register_user", RegisterUser)
	// Создаем mock запрос с отсутствующими данными пользователя
	req, _ := http.NewRequest("GET", "/register_user", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	// Проверяем, что запрос возвращает ошибку
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
