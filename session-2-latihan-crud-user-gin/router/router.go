package router

import (
	"solution1/session-2-latihan-crud-user-gin/handler"
	"solution1/session-2-latihan-crud-user-gin/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	usersPublicEndpoint := r.Group("/users")
	usersPublicEndpoint.GET("/:id", handler.GetUser)
	usersPublicEndpoint.GET("/", handler.GetAllUsers)

	usersPrivateEndpoint := r.Group("/users")
	usersPrivateEndpoint.Use(middleware.AuthMiddleware())
	usersPrivateEndpoint.POST("/", handler.CreateUser)
	usersPrivateEndpoint.PUT("/:id", handler.UpdateUser)
	usersPrivateEndpoint.DELETE("/:id", handler.DeleteUser)

}
