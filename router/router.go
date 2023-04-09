package router

import (
	"eleventh-learn/controllers"
	"eleventh-learn/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)
		userRouter.POST("/login", controllers.UserLogin)
	}

	productRouter := r.Group("/products")
	{
		productRouter.Use(middlewares.Authentication())
		productRouter.POST("/", controllers.CreateProduct)
		productRouter.GET("/", controllers.GetAllProducts)
		productRouter.GET("/:productId", controllers.GetProductByID)
		productRouter.PUT("/:productId", middlewares.ProductAuth(), controllers.UpdateProduct)
		productRouter.DELETE("/:productId", middlewares.ProductAuth(), controllers.DeleteProduct)
	}

	return r
}
