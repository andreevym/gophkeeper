package storage

import (
	"context"
)

type User struct {
	ID       uint64 `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

//go:generate mockgen -source=user.go -destination=./mock/user.go -package=mock
type UserStorage interface {
	GetUser(ctx context.Context, id uint64) (User, error)
	GetUserByLogin(ctx context.Context, login string) (User, error)
	CreateUser(ctx context.Context, user User) (User, error)
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id uint64) error
}
