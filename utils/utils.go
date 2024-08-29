package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"main/config"
	"main/models"
	"time"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
func HashToken(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password[0:56]), 10)
	CheckError(err)
	return string(bytes)
}

func CheckTokenHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password[0:56]))
	return err == nil
}
func ExtractClaims(tokenStr string) (jwt.MapClaims, bool) {
	Secret := []byte(config.SecretKey)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return Secret, nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

// Функция для создания Access токена
func CreateAccessToken(userID string, sessionID int, ipaddr string) string {
	claims := &models.Claims{
		Uuid:      userID,
		Sessionid: sessionID,
		Ip:        ipaddr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)), // Access токен действует 10 минут
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	accessToken, err := token.SignedString([]byte(config.SecretKey)) //access
	CheckError(err)

	return accessToken
}

// Функция для создания Refresh токена
func CreateRefreshToken(userID string, sessionID int, ipaddr string) string {
	claims := &models.Claims{
		Uuid:      userID,
		Sessionid: sessionID,
		Ip:        ipaddr,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(14 * 24 * time.Hour)), // Refresh токен действует 14 дней
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	refreshToken, err := token.SignedString([]byte(config.SecretKey)) //refresh
	CheckError(err)
	return refreshToken
}
