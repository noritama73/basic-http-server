package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/finatext/academia-basic-http-server/internal/domain"
	"github.com/finatext/academia-basic-http-server/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// Common errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
)

// UserUseCase handles user-related business logic
type UserUseCase struct {
	userRepo   domain.UserRepository
	jwtManager *util.JWTManager
}

// NewUserUseCase creates a new user use case
func NewUserUseCase(userRepo domain.UserRepository, jwtManager *util.JWTManager) *UserUseCase {
	return &UserUseCase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// RegisterRequest represents the data needed to register a new user
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// RegisterResponse represents the response after user registration
type RegisterResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Register creates a new user
func (uc *UserUseCase) Register(req RegisterRequest) (*RegisterResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, ErrInvalidInput
	}

	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	// Check if username already exists
	_, err := uc.userRepo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("username already taken")
	}

	// Generate a unique ID
	id, err := generateID()
	if err != nil {
		return nil, err
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:        id,
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Store the user
	if err := uc.userRepo.Store(user); err != nil {
		return nil, err
	}

	return &RegisterResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the response after user login
type LoginResponse struct {
	Token string `json:"token"`
}

// Login authenticates a user and returns a JWT token
func (uc *UserUseCase) Login(req LoginRequest) (*LoginResponse, error) {
	// Find the user
	user, err := uc.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := uc.jwtManager.Generate(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
	}, nil
}

// GetUserRequest represents the data needed to get user information
type GetUserRequest struct {
	ID string `json:"id"`
}

// GetUserResponse represents the response for user information
type GetUserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUser retrieves user information
func (uc *UserUseCase) GetUser(req GetUserRequest) (*GetUserResponse, error) {
	user, err := uc.userRepo.FindByID(req.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Helper function to generate a random ID
func generateID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
