package repository

import (
	"context"
	"go-crud-oapi/internal/model"
)

type UserRepoInterface interface {
	Create(ctx context.Context, user *model.User) error
	ListAllUsers(ctx context.Context) ([]model.User, error)
	GetUserById(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, id uint, user *model.User) error
	DeleteUser(ctx context.Context, id uint) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
}
