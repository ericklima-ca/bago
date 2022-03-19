package controllers

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ericklima-ca/bago/models"
	cachingservice "github.com/ericklima-ca/bago/services/caching_service"
	mailingservice "github.com/ericklima-ca/bago/services/mailing_service"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	DB *gorm.DB
}

type loginPayload struct {
	ID       uint   `json:"id,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}
type bagoResponse struct {
	Ok   bool  `json:"ok"`
	Body gin.H `json:"body,omitempty"`
}

func (a *AuthController) Login(c *gin.Context) {
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
	var user models.User

	if result := a.DB.First(&user, "id = ?", loginPayload.ID); result.RowsAffected != 0 {
		ok := user.TryAuthenticate(loginPayload.Password)
		if !ok && !user.Active {
			c.JSON(http.StatusNotImplemented, bagoResponse{
				Ok: false,
				Body: gin.H{
					"error": "user not active",
				},
			})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, bagoResponse{
			Ok: false,
			Body: gin.H{
				"error": "incorrect password or login",
			},
		})
		return
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

	c.JSON(http.StatusOK, br)
}

func (a *AuthController) Signup(c *gin.Context) {
	var userFormData models.UserFormData
	if err := c.ShouldBindJSON(&userFormData); err != nil {
		c.JSON(http.StatusBadRequest, bagoResponse{
			Ok: false,
			Body: gin.H{
				"error": "data does not match",
			},
		})
		return
	}
	if result := a.DB.Create(userFormData.GetUser()); result.Error != nil {
		c.JSON(http.StatusUnauthorized, bagoResponse{
			Ok: false,
			Body: gin.H{
				"error": "user not created",
			},
		})
		return
	}
	// 115.575µs with goroutine
	// 901.311µs without goroutine
	done := make(chan bool)
	go func() {
		cachingservice.SetToken("signup", userFormData.ID)
		mailingservice.SendConfirmationEmail(userFormData.ID, userFormData.Name, userFormData.Email, c.Request.Host)
		done <- true
	}()
	c.JSON(http.StatusCreated, bagoResponse{
		Ok: true,
		Body: gin.H{
			"msg": "user created",
		},
	})
	<-done
	return
}

func (a *AuthController) Recovery(c *gin.Context) {
	var payloadRecovery struct {
		ID       uint   `json:"id" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var user models.User
	if err := c.ShouldBindJSON(&payloadRecovery); err != nil {
		c.JSON(http.StatusBadRequest, bagoResponse{
			Ok: false,
			Body: gin.H{
				"error": "data does not match",
			},
		})
		return
	}

	if result := a.DB.First(&user, "id = ?", payloadRecovery.ID); result.RowsAffected != 0 {
		cachingservice.SetToken("recovery", user.ID, payloadRecovery.Password)
		mailingservice.SendRecoveryEmail(user.ID, user.Name, user.Email, c.Request.Host)
	}

	c.JSON(http.StatusAccepted, bagoResponse{
		Ok: true,
		Body: gin.H{
			"msg": "confirm the changes in your email",
		},
	})
	return
}

func (a *AuthController) Verify(c *gin.Context) {
	userIdParam, _ := strconv.Atoi(c.Param("id"))
	tokenParam := c.Param("token")
	action := c.Param("action")
	switch action {
	case "signup":
		cachedToken := cachingservice.GetToken("signup", uint(userIdParam))
		if tokenParam == cachedToken {
			var user models.User
			a.DB.First(&user, "id = ?", userIdParam).Update("active", true)
			c.JSON(http.StatusOK, bagoResponse{
				Ok: true,
				Body: gin.H{
					"msg": "user has been activated",
				},
			})
			cachingservice.DeleteToken("signup", uint(userIdParam))
			return
		} else {
			c.JSON(http.StatusBadRequest, bagoResponse{
				Ok: false,
				Body: gin.H{
					"error": "failed to activate user or token expirated",
				},
			})
			return
		}
	case "recovery":
		cachedResult := cachingservice.GetToken("recovery", uint(userIdParam))
		tokenPassList := strings.SplitN(cachedResult, ":", 2)
		token, pass := tokenPassList[0], tokenPassList[1]
		if tokenParam == token {
			var user models.User
			a.DB.First(&user, "id = ?", userIdParam).Update("hashed_password", pass)
			c.JSON(http.StatusOK, bagoResponse{
				Ok: true,
				Body: gin.H{
					"msg": "password updated",
				},
			})
			cachingservice.DeleteToken("recovery", uint(userIdParam))
			return
		} else {
			c.JSON(http.StatusBadRequest, bagoResponse{
				Ok: false,
				Body: gin.H{
					"error": "failed to update password or token expirated",
				},
			})
			return
		}

	}
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
