package jwt

import (
	"context"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// UsernameKey is the context key for storing authenticated username
	UsernameKey contextKey = "username"
)

// JWTUnaryInterceptor creates a production-ready JWT validation interceptor
func JWTUnaryInterceptor() grpc.UnaryServerInterceptor {
	// Initialize JWT key from environment
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable is required")
	}

	jwtKey := &JwtKey{
		SecretKey: []byte(secretKey),
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("Invoking method: " + info.FullMethod)

		// Skip authentication for certain methods (like login, register)
		if isPublicMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		// Extract JWT token from metadata
		token, err := extractTokenFromMetadata(ctx)
		if err != nil {
			return nil, err
		}

		// Validate JWT token
		username, err := jwtKey.ValidateToken(token)
		if err != nil {
			log.Printf("JWT validation failed: %v", err)
			return nil, status.Error(codes.Unauthenticated, "Invalid or expired token")
		}

		// Add username to context for downstream handlers
		ctx = context.WithValue(ctx, UsernameKey, username)

		return handler(ctx, req)
	}
}

// extractTokenFromMetadata extracts and validates JWT token from gRPC metadata
func extractTokenFromMetadata(ctx context.Context) (string, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "Missing metadata")
	}

	// Check for Authorization header (Bearer token)
	if authHeaders := meta.Get("authorization"); len(authHeaders) > 0 {
		authHeader := authHeaders[0]
		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			return strings.TrimPrefix(authHeader, authHeader[7:]), nil
		}
	}

	// Fallback to jwt key for backward compatibility
	if jwtHeaders := meta.Get("jwt"); len(jwtHeaders) == 1 {
		token := strings.TrimSpace(jwtHeaders[0])
		if token == "" {
			return "", status.Error(codes.Unauthenticated, "Empty JWT token")
		}
		return token, nil
	}

	return "", status.Error(codes.Unauthenticated, "Missing or invalid authorization token")
}

// isPublicMethod determines if a method should skip JWT validation
func isPublicMethod(method string) bool {
	publicMethods := []string{
		"/alexchatapp.AuthService/Login",
		"/alexchatapp.AuthService/Register",
	}

	for _, publicMethod := range publicMethods {
		if method == publicMethod {
			return true
		}
	}
	return false
}

// GetUsernameFromContext extracts the authenticated username from context
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}
