package main

import (
	"context"
	"fmt"
	"log"
	"net"
	grpcHandler "solution1/session-11-introduction-redis/handler/grpc"
	pb "solution1/session-11-introduction-redis/proto/user_service/v1"
	"solution1/session-11-introduction-redis/repository/postgres_gorm"
	"solution1/session-11-introduction-redis/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 60*time.Second).Err()
	if err != nil {
		panic(err)
	}

	// setup gorm connection
	dsn := "postgresql://postgres:postgres@localhost:5432/postgres"
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatalln(err)
	}
	// setup service

	// uncomment to use postgres gorm
	userRepo := postgres_gorm.NewUserRepository(gormDB)
	userService := service.NewUserService(userRepo, rdb)
	userHandler := grpcHandler.NewUserHandler(userService)

	// Run the grpc server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Println("Running grpc server in port :50051")
		_ = grpcServer.Serve(lis)
	}()
	time.Sleep(1 * time.Second)

	// Run the grpc gateway
	conn, err := grpc.NewClient(
		"0.0.0.0:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	gwmux := runtime.NewServeMux()
	if err = pb.RegisterUserServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	// dengan GIN
	gwServer := gin.Default()
	gwServer.Group("v1/*{grpc_gateway}").Any("", gin.WrapH(gwmux))
	log.Println("Running grpc gateway server in port :8080")
	_ = gwServer.Run()

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist

	val, err = rdb.Get(ctx, "key-saya").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key-saya", val)
}
