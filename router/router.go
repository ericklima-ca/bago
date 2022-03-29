package router

import (
	"os"

	"github.com/ericklima-ca/bago/router/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	AuthController  AuthHandler
	OrderController OrderHandler
}

type OrderHandler interface {
	CreateMany(*gin.Context)
	GetAll(*gin.Context)
	Delete(*gin.Context)
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
	ordersGroup.Use(middlewares.AuthGuard(os.Getenv("JWT_SECRET")))
	{
		ordersGroup.POST("/", r.OrderController.CreateMany)
		ordersGroup.GET("/", r.OrderController.GetAll)
		ordersGroup.DELETE("/:id", r.OrderController.Delete)
	}

	return routerEngine
}
