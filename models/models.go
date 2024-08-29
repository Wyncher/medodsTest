package models

import "github.com/golang-jwt/jwt/v5"

type JsonResponce struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	Uuid    string `json:"uuid"`
}

// Структура регистрации пользователя
type User struct {
	Username string `form:"username",json:"username"`
	Email    string `form:"email",json:"email"`
}

// Структура входящего запроса на получение токена
type GetTokenStruct struct {
	Uuid string `form:"uuid",json:"uuid"`
}

// Структура входящего запроса на обновление токена
type RefreshTokenStruct struct {
	Token string `form:"token",json:"token"`
	Ip    string `form:"ip",json:"ip"`
}

// Структура для хранения данных токена
type Claims struct {
	Uuid      string `json:"uuid"`
	Sessionid int    `json:"sessionid"`
	Ip        string `json:"ip"`
	jwt.RegisteredClaims
}
