package main

import (
	"log"
	"solution1/session-5-validator/entity"
	"solution1/session-5-validator/handler"
	slice "solution1/session-5-validator/repository/slice"
	"solution1/session-5-validator/router"
	"solution1/session-5-validator/service"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// setup service
	var mockUserDBInSlice []entity.User
	userRepo := slice.NewUserRepository(mockUserDBInSlice)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Routes
	router.SetupRouter(r, userHandler)

	// Run the server
	log.Println("Running server on port 8080")
	r.Run(":8080")
}
