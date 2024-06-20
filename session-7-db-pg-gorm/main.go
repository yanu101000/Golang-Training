package main

import (
	"context"
	"log"
	"solution1/session-7-db-pg-gorm/handler"
	"solution1/session-7-db-pg-gorm/repository/postgres_gorm"
	"solution1/session-7-db-pg-gorm/repository/postgres_pgx"
	"solution1/session-7-db-pg-gorm/router"
	"solution1/session-7-db-pg-gorm/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// setup database connection
	//dsn := "postgresql://postgres:postgres@localhost:5432/postgres"
	// setup pgx connection
	//pgxPool, err := connectDB(dsn)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// setup gorm connectoin
	dsn := "postgresql://postgres:postgres@localhost:5432/postgres"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}
	// setup service

	// slice db is disabled. uncomment to enabled
	// var mockUserDBInSlice []entity.User
	// _ = slice.NewUserRepository(mockUserDBInSlice)

	// uncomment to use postgres pgx
	// userRepo := postgres_pgx.NewUserRepository(pgxPool)

	// uncomment to use postgres gorm
	userRepo := postgres_gorm.NewUserRepository(gormDB)

	// service and handler declaration
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Routes
	router.SetupRouter(r, userHandler)

	// Run the server
	log.Println("Running server on port 8080")
	r.Run(":8080")
}

func connectDB(dbURL string) (postgres_pgx.PgxPoolIface, error) {
	return pgxpool.New(context.Background(), dbURL)
}
