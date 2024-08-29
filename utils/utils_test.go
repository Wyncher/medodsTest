package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"main/config"
	"testing"
	"time"
)

// Тест функции ошибки
func TestCheckError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Функция не паникует")
		}
	}()

	// вызываем функцию с ошибкой
	CheckError(errors.New("Ошибка"))
}

// Тест хэширования токена и проверки
func TestHashTokenAndCheckTokenHash(t *testing.T) {
	password := "оль13оо2о2Пароль13оо2о2Пароль13оо2о2Пароль13оо2о2Пароль13оо2о2Пароль13оо2о2Пароль13оо2о2Пароль13оо2о2Пароль13оо2о2"
	hash := HashToken(password)

	if !CheckTokenHash(password, hash) {
		t.Errorf("Пароль не подходит")
	}

	// тестируем с неверным паролем
	if CheckTokenHash("wrongpasswordwrongpasswordwrongpasswordwrongpasswordwrongpasswordwrongpasswordwrongpasswordwrongpassword", hash) {
		t.Errorf("Пароль не правильный")
	}
}

// Тестирование извлечения данных с токена
func TestExtractClaims(t *testing.T) {
	// Создаем валидный токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": "123",
		"exp":     time.Now().Add(10 * time.Minute).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(config.SecretKey))

	claims, valid := ExtractClaims(tokenString)
	if !valid {
		t.Errorf("Данные с токена не извлечены")
	}
	if claims["user_id"] != "123" {
		t.Errorf("Ождиаемый user_id = '123', получен '%v'", claims["user_id"])
	}

	// Тестируем с невалидным токеном
	claims, valid = ExtractClaims("invalid_token")
	if valid {
		t.Errorf("Токен недействителен")
	}
}

// Тестирование создания токенов
func TestCreateTokens(t *testing.T) {
	userID := "99a79549-b0bc-49f4-8f36-2c5be1e5e805"
	sessionID := 1
	ipaddr := "127.0.0.1"

	accessToken := CreateAccessToken(userID, sessionID, ipaddr)
	refreshToken := CreateRefreshToken(userID, sessionID, ipaddr)

	if accessToken == "" {
		t.Errorf("Токен доступа не может быть пустым")
	}
	if refreshToken == "" {
		t.Errorf("Токен обновления не может быть пустым")
	}

	// Проверяем на валидность и полноту данных
	accessClaims, accessValid := ExtractClaims(accessToken)
	if !accessValid {
		t.Errorf("Токен доступа не действителен")
	}
	if accessClaims["uuid"].(string) != userID {
		t.Errorf("Ожидаемый Uuid в токене доступа = '%v', получен '%v'", userID, accessClaims["Uuid"])
	}

	refreshClaims, refreshValid := ExtractClaims(refreshToken)
	if !refreshValid {
		t.Errorf("Токен обновления не действителен")
	}
	if refreshClaims["uuid"].(string) != userID {
		t.Errorf("Ожидаемый Uuid в токене доступа = '%v', получен '%v'", userID, refreshClaims["Uuid"])
	}
}
