package domain

import (
	"time"
)

// User represents a user entity in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Password is not included in JSON responses
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	FindByID(id string) (*User, error)
	FindByUsername(username string) (*User, error)
	Store(user *User) error
	Update(user *User) error
	Delete(id string) error
}
