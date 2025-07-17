package repository

import (
	"context"
	"go-crud-oapi/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepoInterface {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepo) GetUserById(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, id uint, user *model.User) error {
	return r.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *UserRepo) DeleteUser(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	if r.DB.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *UserRepo) ListAllUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.DB.WithContext(ctx).Find(&users).Error
	return users, err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Email not found, it's okay
		}
		return nil, err // Actual DB error
	}
	return &user, nil // Email found
}
