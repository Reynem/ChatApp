package alexchatapp

import (
	"alexchatapp/src/data"
	pb "alexchatapp/src/proto/auth"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"

	"google.golang.org/grpc"
)

func Server() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	dsn := os.Getenv("POSTGRES_CONNECTION")
	db, err := data.Initialize(dsn)
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// Create repositories
	chat_repo := data.NewChatRepository(db)
	auth_repo := data.NewAuthRepository(db)

	// Create authentication server
	authServer := NewAuthServer(chat_repo, auth_repo)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	// Start server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error creating listener: %v", err)
	}

	log.Println("Server started on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
