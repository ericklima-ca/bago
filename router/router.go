package router

import (
	"os"

	"github.com/ericklima-ca/bago/controllers"
	"github.com/gin-gonic/gin"
)

func LoadRoutes() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	r := gin.Default()

	auther := r.Group("/api/auth")
	auther.POST("/login", controllers.Auth.Login)
	auther.POST("/signup", controllers.Auth.Signup)
	auther.HEAD("/recovery/:id", controllers.Auth.Recovery)

	return r
}
