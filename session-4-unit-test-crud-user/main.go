package main

import (
	"log"
	"solution1/session-4-unit-test-crud-user/entity"
	"solution1/session-4-unit-test-crud-user/handler"
	slice "solution1/session-4-unit-test-crud-user/repository/slice"
	"solution1/session-4-unit-test-crud-user/router"
	"solution1/session-4-unit-test-crud-user/service"

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
