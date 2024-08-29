package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apis "main/apis"
	"main/config"
)

type Handler struct{}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/")
	v1.GET("/register", apis.RegisterUser) //принимает username и email и возвращает uuid
	v1.GET("/getToken", apis.GetToken)     //принимает uuid и отдает пару токенов
	v1.GET("/refresh", apis.RefreshToken)  //принимает refresh и отдает пару токенов
	r.Run(fmt.Sprintf(":%v", config.ServerPort))
}
