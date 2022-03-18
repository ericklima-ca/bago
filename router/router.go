package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	AuthController AuthHandler
}

type AuthHandler interface {
	Login(*gin.Context)
	Signup(*gin.Context)
	Recovery(*gin.Context)
	Verify(*gin.Context)
}

func (r *Router) LoadRoutes() *gin.Engine {
	routerEngine := gin.Default()

	authGroup := routerEngine.Group("/api/auth")
	{
		authGroup.POST("/login", r.AuthController.Login)
		authGroup.POST("/signup", r.AuthController.Signup)
		authGroup.POST("/recovery/", r.AuthController.Recovery)
		authGroup.GET("/verify/:action/:id/:token", r.AuthController.Verify)

	}

	return routerEngine
}
