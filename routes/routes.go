package routes

import (
	"gitlab.com/ployMatsuri/go-backend/controllers"

	"gitlab.com/ployMatsuri/go-backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	product := r.Group("/products")
	product.Use(middleware.AuthMiddleware())
	{
		product.GET("/", controllers.GetProducts)
	}

	return r
}
