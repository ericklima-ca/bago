package router

import (
	"github.com/ericklima-ca/bago/router/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	AuthController  AuthHandler
	OrderController OrderHandler
}

type OrderHandler interface {
	Create(*gin.Context)
}

type AuthHandler interface {
	Login(*gin.Context)
	Signup(*gin.Context)
	Recovery(*gin.Context)
	Verify(*gin.Context)
}

func (r *Router) LoadRoutes() *gin.Engine {
	routerEngine := gin.Default()
	routerEngine.Use(cors.Default())

	authGroup := routerEngine.Group("/api/auth")
	{
		authGroup.POST("/login", r.AuthController.Login)
		authGroup.POST("/signup", r.AuthController.Signup)
		authGroup.POST("/recovery/", r.AuthController.Recovery)
		authGroup.GET("/verify/:action/:id/:token", r.AuthController.Verify)

	}
	ordersGroup := routerEngine.Group("/api/orders")
	{
		ordersGroup.Use(middlewares.AuthGuard())
		ordersGroup.POST("/create", r.OrderController.Create)
	}

	return routerEngine
}
