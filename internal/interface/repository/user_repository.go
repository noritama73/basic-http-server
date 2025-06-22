package repository

import (
	"errors"
	"sync"

	"github.com/finatext/academia-basic-http-server/internal/domain"
)

// Common errors
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

// InMemoryUserRepository is an in-memory implementation of UserRepository
type InMemoryUserRepository struct {
	mutex sync.RWMutex
	users map[string]*domain.User // Map of user ID to user
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// FindByID finds a user by ID
func (r *InMemoryUserRepository) FindByID(id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// FindByUsername finds a user by username
func (r *InMemoryUserRepository) FindByUsername(username string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// Store adds a new user
func (r *InMemoryUserRepository) Store(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if user with same ID already exists
	if _, exists := r.users[user.ID]; exists {
		return ErrUserExists
	}

	// Check if username is already taken
	for _, existingUser := range r.users {
		if existingUser.Username == user.Username {
			return ErrUserExists
		}
	}

	// Store the user
	r.users[user.ID] = user
	return nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return ErrUserNotFound
	}

	// Check if username is already taken by another user
	for id, existingUser := range r.users {
		if existingUser.Username == user.Username && id != user.ID {
			return ErrUserExists
		}
	}

	r.users[user.ID] = user
	return nil
}

// Delete removes a user
func (r *InMemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}
