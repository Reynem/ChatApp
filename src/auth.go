package alexchatapp

import (
	"context"

	"alexchatapp/src/data"
	pb "alexchatapp/src/proto/auth"
	"alexchatapp/src/utils"
)

// AuthServer implements AuthService from proto file
type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	chat_repo *data.ChatRepository
	auth_repo *data.AuthRepository
}

// NewAuthServer creates a new authentication server instance
func NewAuthServer(chat_repo *data.ChatRepository, auth_repo *data.AuthRepository) *AuthServer {
	return &AuthServer{
		chat_repo: chat_repo,
		auth_repo: auth_repo,
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
	_, err := s.auth_repo.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		return &pb.Response{
			Success:   false,
			ErrorText: err.Error(),
		}, nil
	}

	// Generate token
	// token, err := generateToken()
	// if err != nil {
	// 	return &pb.Response{
	// 		Success:   false,
	// 		ErrorText: "Error generating token",
	// 	}, nil
	// }

	return &pb.Response{
		Success: true,
		Token:   "token",
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

	_, err := s.auth_repo.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		return &pb.Response{
			Success:   false,
			ErrorText: err.Error(),
		}, nil
	}

	// Generate token
	// token, err := generateToken()
	// if err != nil {
	// 	return &pb.Response{
	// 		Success:   false,
	// 		ErrorText: "Error generating token",
	// 	}, nil
	// }

	return &pb.Response{
		Success: true,
		Token:   "token",
	}, nil
}
