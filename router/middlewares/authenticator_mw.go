package middlewares

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ericklima-ca/bago/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"ok": false,
				"body": gin.H{
					"error": "user not authenticated",
				},
			})
		} else {
			t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			if rawPayload, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
				var user models.User
				bytesPayload, _ := json.Marshal(rawPayload)
				json.Unmarshal(bytesPayload, &user)
				c.Set("user", user)
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"ok": false,
					"body": gin.H{
						"error": "invalid jwt token",
					},
				})
			}

		}
	}
}
