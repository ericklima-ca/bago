package controllers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ericklima-ca/bago/models"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	DB *gorm.DB
}

type loginPayload struct {
	Login    string `json:"login,omitempty" binding:"required"`
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
	_id, _ := strconv.Atoi(loginPayload.Login)

	if result := a.DB.First(&user, "id = ?", _id); result.RowsAffected != 0 {
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
	// <-- TBD
	// c.SetCookie("Authorization", "Bearer: "+br.Body["token"].(string), int(time.Now().Add(time.Hour*8).Unix()), "/login", c.Request.Host, true, true)
	// c.Header("Authorization", "Bearer: "+br.Body["token"].(string))
	// ->
	c.JSON(http.StatusOK, br)
}

func (a *AuthController) Recovery(c *gin.Context) {
	var user models.User
	userIdInt, _ := strconv.Atoi(c.Param("id"))
	if result := a.DB.First(&user, "id = ?", userIdInt); result.Error != nil {
		return
	}

	tokenRecovery := models.TokenRecovery{
		TokenBase: models.TokenBase{
			UserID: user.ID,
		},
	}
	a.DB.Create(&tokenRecovery)

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
	token := models.TokenSignup{
		TokenBase: models.TokenBase{
			UserID: userFormData.ID,
		},
	}
	a.DB.Create(&token)
	c.JSON(http.StatusCreated, bagoResponse{
		Ok: true,
		Body: gin.H{
			"msg": "user created",
		},
	})
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
