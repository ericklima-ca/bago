package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/ericklima-ca/bago/models"
	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
)

// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		panic(err)
// 	}
// }

var (
	Auth = authenticator{
		Login: login,
	}
)

type authenticator struct {
	Login gin.HandlerFunc
}
type loginPayload struct {
	Login    string `json:"login,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}
type bagoResponse struct {
	Ok   bool  `json:"ok"`
	Body gin.H `json:"body,omitempty"`
}

func login(c *gin.Context) {
	var loginPayload loginPayload
	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		c.JSON(http.StatusBadRequest, bagoResponse{
			Ok: false,
			Body: gin.H{
				"error": "invalid payload",
			},
		})
		return
	}
	user, ok := models.TryAuthenticate(loginPayload)
	if !ok {
		c.JSON(http.StatusUnauthorized, bagoResponse{
			Ok: false,
			Body: gin.H{
				"error": "incorrect password or login",
			},
		})
	}

	br := bagoResponse{
		Ok: true,
		Body: gin.H{
			"user":  user,
			"token": "",
		},
	}
	if err := setAuthToken(&br); err != nil {
		panic(err)
	}
	// <-- TBD
	// c.SetCookie("Authorization", "Bearer: "+br.Body["token"].(string), int(time.Now().Add(time.Hour*8).Unix()), "/login", c.Request.Host, true, true)
	// c.Header("Authorization", "Bearer: "+br.Body["token"].(string))
	// ->
	c.JSON(http.StatusOK, br)
}

func setAuthToken(br *bagoResponse) error {
	SECRET := os.Getenv("JWT_SECRET")

	claim := jwt.MapClaims{
		"user": br.Body["user"],
		"token": jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 8).Unix(),
		},
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	token, err := tk.SignedString([]byte(SECRET))
	if err != nil {
		panic(err)
	}
	br.Body["token"] = token
	return err
}
