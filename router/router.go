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
}

func (r *Router) LoadRoutes() *gin.Engine {
	routerEngine := gin.Default()

	authGroup := routerEngine.Group("/api/auth")
	authGroup.POST("/login", r.AuthController.Login)
	authGroup.POST("/signup", r.AuthController.Signup)
	authGroup.HEAD("/recovery/:id", r.AuthController.Recovery)

	return routerEngine
}
