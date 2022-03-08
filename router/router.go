package router

import (
	"github.com/ericklima-ca/bago/controllers"
	"github.com/gin-gonic/gin"
)

func LoadRoutes() *gin.Engine {
	r := gin.Default()

	auther := r.Group("/api/auth")
	auther.POST("/login", controllers.Auth.Login)

	return r
}
