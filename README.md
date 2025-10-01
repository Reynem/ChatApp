# AlexChatApp

gRPC-based chat application with JWT authentication and user profiles.

## Features

- **Authentication**: User registration and login with JWT tokens
- **Profiles**: User profile management with bio, avatar, and status
- **Chat**: Real-time messaging (gRPC) (CURRENTLY NOT REALISED)
- **Security**: JWT-based authentication with interceptors

## Quick Start

1. **Setup environment**
   ```bash
   cp .env.example .env
   # Configure POSTGRES_CONNECTION and SECRET_KEY
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run server**
   ```bash
   go run main.go
   ```

Server starts on `:50051`

## API Services

### Auth Service
- `Register(username, email, password)` - Create new user
- `Login(username, password)` - Authenticate user

### Profile Service
- `CreateProfile(name, bio, avatar, status)` - Create user profile
- `GetProfile()` - Get current user profile
- `UpdateProfile(...)` - Update profile data
- `UpdateOnlineStatus(last_seen)` - Update activity status

### Chat Service (IN DEVELOPMENT)
- Message sending and receiving

## Testing

```bash
# Run examples
go run src/client/main.go
```

## Tech Stack

- **Go** - Backend language
- **gRPC** - API protocol
- **PostgreSQL** - Database
- **GORM** - ORM
- **JWT** - Authentication
- **Protocol Buffers** - Message serialization

## Project Structure

```
src/
├── auth.go              # Auth service implementation
├── profiles.go          # Profile service implementation  
├── server.go           # gRPC server setup
├── jwt/                # JWT utilities
├── data/               # Database repositories
├── models/             # Data models
├── proto/              # Protocol buffer definitions
└── tests/              # Client examples
```