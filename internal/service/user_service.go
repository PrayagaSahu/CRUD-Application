package service

import (
	"context"
	"go-crud-oapi/internal/model"
	"go-crud-oapi/internal/repository"

	"gorm.io/gorm"
)

type UserServiceInterFace interface {
	Create(ctx context.Context, user *model.User) error
	ListAllUsers(ctx context.Context) ([]model.User, error)
	Get(ctx context.Context, id uint) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, id uint, user *model.User) error
	Delete(ctx context.Context, id uint) error
}

type UserService struct {
	repo repository.UserRepoInterface
}

func NewUserService(repo repository.UserRepoInterface) UserServiceInterFace {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, user *model.User) error {
	return s.repo.Create(ctx, user)
}

func (s *UserService) ListAllUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.ListAllUsers(ctx)
}

func (s *UserService) Get(ctx context.Context, id uint) (*model.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *UserService) Update(ctx context.Context, id uint, user *model.User) error {
	user.ID = id
	return s.repo.UpdateUser(ctx, id, user)
}

func (s *UserService) Delete(ctx context.Context, id uint) error {
	user, err := s.repo.GetUserById(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return gorm.ErrRecordNotFound
	}

	return s.repo.DeleteUser(ctx, id)
}
