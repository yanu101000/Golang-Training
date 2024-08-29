package main

import (
	"log"
	"net"

	"wallet/wallet-service/handler"
	pb "wallet/wallet-service/proto" // Adjust the import path to your proto package
	"wallet/wallet-service/repository"
	"wallet/wallet-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	// Listen on the specified port
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Database connection string
	dsn := "postgresql://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "assignment_6.", // schema name
			SingularTable: false,
		}})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Initialize repository and service
	walletRepo := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	// Set up gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterWalletServiceServer(grpcServer, walletHandler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	log.Printf("gRPC server started at %s", listen.Addr().String())
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
