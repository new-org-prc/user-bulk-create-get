package storage

import (
	"context"
	"fmt"
	"sika/internal/user"
	"sika/pkg/storage/entities"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) user.Repo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) CreateUser(ctx context.Context, u *entities.User) error {
	if err := r.db.WithContext(ctx).Save(&u).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepo) CreateBatchUsers(ctx context.Context, users []entities.User) error {
	if err := r.db.WithContext(ctx).CreateInBatches(users, 10).Error; err != nil {
		return err
	}
	return nil
}
func (r *userRepo) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	result := r.db.WithContext(ctx).Preload("Addresses").First(&user, "id=?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepo) ClearAllUsersDataFromDB() error {
	if err := r.db.Exec("DELETE FROM users").Error; err != nil {
		return fmt.Errorf("failed to clear users table: %w", err)
	}
	return nil
}
