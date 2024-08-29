package apis

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"main/config"
	"main/db"
	"main/models"
	"main/utils"
	"math/rand"
	"net/http"
	"time"
)

var conn = db.Conn()

// роут на обновление токенов
func RefreshToken(c *gin.Context) {
	var token models.RefreshTokenStruct
	if err := c.BindQuery(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var NewToken = ""
	var OldToken = token.Token
	ip := c.ClientIP()
	claims, errror := utils.ExtractClaims(OldToken) //структура с токена refresh
	if !errror {
		return
	}
	ttl := claims["exp"].(float64) //время жизни
	claimIp := claims["ip"].(string)
	uuidClaim := claims["uuid"].(string)       //uuid
	expirationTime := time.Unix(int64(ttl), 0) //сколько ttl

	//Если IP адрес изменен
	if ip != claimIp {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Это не вы!"})
		sendAlert(uuidClaim)
		return
	}
	// Текущая метка времени
	currentTime := time.Now()
	if currentTime.After(expirationTime) {
		//Токен просрочен
		c.JSON(http.StatusOK, gin.H{"error": "Время жизни токена закончилось. Войдите снова"})
	} else {
		//Токен не просрочен
		session_ID := int(claims["sessionid"].(float64))
		//Поиск bcrypt токена в БД
		rows, err := conn.Query(`SELECT searchrefresh($1,$2);`, ip, session_ID)
		utils.CheckError(err)
		for rows.Next() {
			//Чтение токена
			rows.Scan(&NewToken)
		}
		//Если токен найден в БД
		if NewToken != "" {
			//Если токен bcrypt в БД совпадает с текущим токеном от пользователя
			if utils.CheckTokenHash(OldToken, NewToken) {
				fmt.Println("bingo")

			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Это не вы!"})
				sendAlert(uuidClaim)
				return
			}

		} else {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Токен не найден в БД"})
			sendAlert(uuidClaim)
			return
		}

	}
	var sessionId = rand.Intn(config.RandLength)
	accessToken := utils.CreateAccessToken(uuidClaim, sessionId, ip)
	//Обновление refresh токена в БД
	refreshToken := utils.CreateRefreshToken(uuidClaim, sessionId, ip)
	bcryptToken := utils.HashToken(refreshToken)
	_, err := conn.Query(`SELECT updaterefreshtoken($1,$2,$3,$4);`, uuidClaim, bcryptToken, ip, sessionId)
	utils.CheckError(err)
	response, _ := json.Marshal(models.JsonResponce{Access: accessToken, Refresh: OldToken, Uuid: uuidClaim})
	c.Data(http.StatusOK, gin.MIMEJSON, response)
	return
}

// роут на получение токенов
func GetToken(c *gin.Context) {
	var uuID models.GetTokenStruct
	if err := c.BindQuery(&uuID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var uuid = uuID.Uuid
	var sessionID = rand.Intn(config.RandLength)
	ip := c.ClientIP()
	if ip == "" {
		ip = "noip"
	}
	accessToken := utils.CreateAccessToken(uuid, sessionID, ip)
	refreshToken := utils.CreateRefreshToken(uuid, sessionID, ip)
	bcryptToken := utils.HashToken(refreshToken)
	_, err := conn.Query(`SELECT setrefreshtoken($1,$2,$3,$4);`, uuid, bcryptToken, ip, sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, _ := json.Marshal(models.JsonResponce{Access: accessToken, Refresh: refreshToken, Uuid: uuid})
	c.Data(http.StatusOK, gin.MIMEJSON, response)
}

// роут на регистрацию пользователя
func RegisterUser(c *gin.Context) {
	var user models.User
	var uuID = uuid.New().String()
	if err := c.BindQuery(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := conn.Query(`SELECT register($1, $2, $3);`, user.Username,
		user.Email, uuID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, gin.MIMEJSON, []byte(uuID))

}

// Функция отправки предупреждения на почту
func sendAlert(uuid string) {
	var email string
	rows, err := conn.Query(`SELECT searchemail($1);`, uuid)
	utils.CheckError(err)
	for rows.Next() {
		rows.Scan(&email)
	}
	//тут может быть отправка уведомления на почту

}
