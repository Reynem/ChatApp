package alexchatapp

import (
	"alexchatapp/src/data"
	"alexchatapp/src/jwt"
	pba "alexchatapp/src/proto/auth"
	pbp "alexchatapp/src/proto/profiles"
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

	// Getting sercret variable
	secret_key := os.Getenv("SECRET_KEY")
	jwt_key := jwt.JwtKey{SecretKey: []byte(secret_key)}

	// Create repositories
	chat_repo := data.NewChatRepository(db)
	auth_repo := data.NewUsersRepository(db)
	profile_repo := data.NewProfilesRepository(db)

	// Create authentication server
	authServer := NewAuthServer(chat_repo, auth_repo, profile_repo, &jwt_key)
	profileServer := NewProfilesServer(profile_repo)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(jwt.JWTUnaryInterceptor()),
	)

	pba.RegisterAuthServiceServer(grpcServer, authServer)
	pbp.RegisterProfileServiceServer(grpcServer, profileServer)

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
