package tests

import (
	"context"
	"log"
	"time"

	pba "alexchatapp/src/proto/auth"
	pb "alexchatapp/src/proto/profiles"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProfileClientExample demonstrates profile service client usage
func ProfileClientExample() {
	// Connect to server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection error: %v", err)
	}
	defer conn.Close()

	// First, authenticate to get JWT token
	authClient := pba.NewAuthServiceClient(conn)
	profileClient := pb.NewProfileServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Login to get JWT token
	log.Println("Logging in to get JWT token...")
	loginResp, err := authClient.Login(ctx, &pba.LoginRequest{
		Username: "testuser", // Make sure this user exists
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login error: %v", err)
	}

	if !loginResp.Success {
		log.Fatalf("Login failed: %s", loginResp.ErrorText)
	}

	log.Printf("✅ Login successful! Token received")

	// Create context with JWT token in metadata
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + loginResp.Token,
	})
	authCtx := metadata.NewOutgoingContext(ctx, md)

	// === 1. Регистрация нового профиля ===
	// log.Println("Registering new profile...")
	bio := "Hello, I'm a test user!"
	avatarUrl := "https://example.com/avatar.png"
	status := "active"

	// log.Println("Using the token: " + loginResp.Token)
	// registerResp, err := profileClient.CreateProfile(authCtx, &pb.CreateProfileRequest{
	// 	ProfileName: "TestProfile",
	// 	Bio:         &bio,
	// 	AvatarUrl:   &avatarUrl,
	// 	Status:      &status,
	// })
	// if err != nil {
	// 	log.Fatalf("RegisterProfile error: %v", err)
	// }

	// if registerResp.StatusCode == 200 {
	// 	log.Println("✅ Profile registration successful!")
	// } else {
	// 	log.Printf("❌ Profile registration failed with status: %d", registerResp.StatusCode)
	// }

	// === 2. Попытка регистрации дубликата ===
	log.Println("Trying to register duplicate profile...")
	duplicateResp, err := profileClient.CreateProfile(authCtx, &pb.CreateProfileRequest{
		ProfileName: "AnotherName",
		Bio:         &bio,
		AvatarUrl:   &avatarUrl,
		Status:      &status,
	})
	if err != nil {
		log.Printf("Expected error on duplicate: %v", err)
	} else if duplicateResp.StatusCode == 400 {
		log.Println("✅ Correctly rejected duplicate profile")
	} else {
		log.Printf("❌ Unexpected response for duplicate: status=%d", duplicateResp.StatusCode)
	}

	// === 3. Получение профиля ===
	log.Println("Fetching profile...")
	getResp, err := profileClient.GetProfile(authCtx, &pb.GetProfileRequest{})
	if err != nil {
		log.Fatalf("GetProfile error: %v", err)
	}

	if getResp.Profile != nil {
		log.Printf("✅ Retrieved profile: %s (Bio: %s)", getResp.Profile.ProfileName, *getResp.Profile.Bio)
	} else {
		log.Println("❌ Failed to retrieve profile")
	}

	// === 4. Обновление профиля ===
	log.Println("Updating profile...")
	newBio := "Updated bio!"
	newAvatar := "https://example.com/new-avatar.png"
	newStatus := "busy"
	updatedProfileName := "UpdatedTestProfile"

	updateResp, err := profileClient.UpdateProfile(authCtx, &pb.UpdateProfileRequest{
		ProfileName: &updatedProfileName,
		Bio:         &newBio,
		AvatarUrl:   &newAvatar,
		Status:      &newStatus,
	})
	if err != nil {
		log.Fatalf("UpdateProfile error: %v", err)
	}

	if updateResp.StatusCode == 200 {
		log.Println("✅ Profile updated successfully!")
	} else {
		log.Printf("❌ Profile update failed with status: %d", updateResp.StatusCode)
	}

	// === 5. Обновление онлайн-статуса (LastSeen) ===
	log.Println("Updating online status (LastSeen)...")
	lastSeen := timestamppb.Now()

	onlineResp, err := profileClient.UpdateOnlineStatus(authCtx, &pb.UpdateOnlineStatusRequest{
		LastSeen: lastSeen,
	})
	if err != nil {
		log.Fatalf("UpdateOnlineStatus error: %v", err)
	}

	if onlineResp.StatusCode == 200 {
		log.Println("✅ Online status updated successfully!")
	} else {
		log.Printf("❌ Online status update failed with status: %d", onlineResp.StatusCode)
	}

	// === 6. Получение обновленного профиля ===
	log.Println("Fetching updated profile...")
	finalResp, err := profileClient.GetProfile(authCtx, &pb.GetProfileRequest{})
	if err != nil {
		log.Printf("Final GetProfile error: %v", err)
	} else if finalResp.Profile != nil {
		log.Printf("✅ Final profile state:")
		log.Printf("   Name: %s", finalResp.Profile.ProfileName)
		log.Printf("   Bio: %s", *finalResp.Profile.Bio)
		log.Printf("   Status: %s", *finalResp.Profile.Status)
		log.Printf("   Avatar: %s", *finalResp.Profile.AvatarUrl)
	}

	log.Println("✅ All profile service tests completed!")
}
