package alexchatapp

import (
	"context"
	"log"
	"time"

	"alexchatapp/src/data"
	"alexchatapp/src/jwt"
	"alexchatapp/src/models"
	pb "alexchatapp/src/proto/auth"
	"alexchatapp/src/utils"
)

// AuthServer implements AuthService from proto file
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	chat_repo    *data.ChatRepository
	auth_repo    *data.UsersRepository
	profile_repo *data.ProfilesRepository
	jwtService   *jwt.JwtKey
}

// NewAuthServer creates a new authentication server instance
func NewAuthServer(chat_repo *data.ChatRepository, auth_repo *data.UsersRepository, profile_repo *data.ProfilesRepository, secret_key *jwt.JwtKey) *AuthServer {
	return &AuthServer{
		chat_repo:    chat_repo,
		auth_repo:    auth_repo,
		profile_repo: profile_repo,
		jwtService:   secret_key,
	}
}

// Register registers a new user
func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.Response, error) {
	// Validate input data
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return &pb.Response{
			Success:   false,
			ErrorText: "All fields are required",
		}, nil
	}

	// Validate username
	if err := utils.ValidateUsername(req.Username); err != nil {
		return &pb.Response{
			Success:   false,
			ErrorText: err.Error(),
		}, nil
	}

	// Validate email
	if !utils.ValidateEmail(req.Email) {
		return &pb.Response{
			Success:   false,
			ErrorText: "Invalid email address",
		}, nil
	}

	// Validate password
	if err := utils.ValidatePassword(req.Password); err != nil {
		return &pb.Response{
			Success:   false,
			ErrorText: err.Error(),
		}, nil
	}

	// Add user to database
	user, err := s.auth_repo.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		return &pb.Response{
			Success:   false,
			ErrorText: err.Error(),
		}, nil
	}

	var profile = models.Profile{
		User_id:      user.ID,
		Profile_name: user.UserName,
		Bio:          "",
		Avatar_url:   "",
		Status:       "",
		Last_seen:    time.Now(),
	}

	err = s.profile_repo.CreateProfileByModel(profile)
	if err != nil {
		return &pb.Response{
			Success:   false,
			ErrorText: err.Error(),
		}, nil
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user.UserName, uint64(user.ID))
	if err != nil {
		log.Printf("JWT generation error: %v", err)
		return &pb.Response{
			Success:   false,
			ErrorText: "Error generating token",
		}, nil
	}

	return &pb.Response{
		Success: true,
		Token:   token,
	}, nil
}

// Login performs user authentication
func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.Response, error) {
	// Validate input data
	if req.Username == "" || req.Password == "" {
		return &pb.Response{
			Success:   false,
			ErrorText: "Username and password are required",
		}, nil
	}

	user, err := s.auth_repo.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		return &pb.Response{
			Success:   false,
			ErrorText: err.Error(),
		}, nil
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(user.UserName, uint64(user.ID))
	if err != nil {
		log.Printf("JWT generation error: %v", err)
		return &pb.Response{
			Success:   false,
			ErrorText: "Error generating token",
		}, nil
	}

	return &pb.Response{
		Success: true,
		Token:   token,
	}, nil
}
