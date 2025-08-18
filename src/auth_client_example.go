package alexchatapp

import (
	"context"
	"log"
	"time"

	pb "alexchatapp/src/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClientExample demonstrates authentication client usage
func AuthClientExample() {
	// Connect to server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Registration example
	log.Println("Registering new user...")
	registerResp, err := client.Register(ctx, &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Registration error: %v", err)
	}

	if registerResp.Success {
		log.Printf("Registration successful! Token: %s", registerResp.Token)
	} else {
		log.Printf("Registration error: %s", registerResp.ErrorText)
	}

	// Login example
	log.Println("Logging in...")
	loginResp, err := client.Login(ctx, &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login error: %v", err)
	}

	if loginResp.Success {
		log.Printf("Login successful! Token: %s", loginResp.Token)
	} else {
		log.Printf("Login error: %s", loginResp.ErrorText)
	}
}
