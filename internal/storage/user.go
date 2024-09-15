package storage

import (
	"context"
)

// User represents a user entity with basic information such as ID, login, and password.
type User struct {
	ID       uint64 `json:"id"`       // Unique identifier for the user.
	Login    string `json:"login"`    // User's login name.
	Password string `json:"password"` // User's password, typically hashed.
}

// UserStorage defines the interface for operations on user entities in the storage system.
type UserStorage interface {
	// GetUser retrieves a user by their unique ID.
	// Takes a context.Context and the user's ID (uint64) as parameters.
	// Returns the User and an error if any.
	GetUser(ctx context.Context, id uint64) (User, error)

	// GetUserByLogin retrieves a user by their login name.
	// Takes a context.Context and the user's login (string) as parameters.
	// Returns the User and an error if any.
	GetUserByLogin(ctx context.Context, login string) (User, error)

	// CreateUser inserts a new user into the storage system.
	// Takes a context.Context and a User object as parameters.
	// Returns the created User and an error if any.
	CreateUser(ctx context.Context, user User) (User, error)

	// UpdateUser updates an existing user's information.
	// Takes a context.Context and a User object with updated information as parameters.
	// Returns an error if any.
	UpdateUser(ctx context.Context, user User) error

	// DeleteUser removes a user from the storage system by their ID.
	// Takes a context.Context and the user's ID (uint64) as parameters.
	// Returns an error if any.
	DeleteUser(ctx context.Context, id uint64) error
}
