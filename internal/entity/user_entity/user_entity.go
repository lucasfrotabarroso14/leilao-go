package user_entity

import (
	"context"
	"leilao-go/internal/internal_error"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserRepositoryInterface interface {
	FindUserById(ctx context.Context, id string) (*User, *internal_error.InternalError)
}
