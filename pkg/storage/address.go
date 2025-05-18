package storage

import (
	"context"
	"fmt"
	"sika/internal/address"
	"sika/pkg/storage/entities"

	"gorm.io/gorm"
)

type addressRepo struct {
	db *gorm.DB
}

func NewAddressRepo(db *gorm.DB) address.Repo {
	return &addressRepo{
		db: db,
	}
}

func (r *addressRepo) CreateAddress(ctx context.Context, a *entities.Address) error {
	if err := r.db.WithContext(ctx).Save(&a).Error; err != nil {
		return err
	}

	return nil
}

func (r *addressRepo) CreateBatchAddresses(ctx context.Context, adds []entities.Address) error {
	if err := r.db.WithContext(ctx).CreateInBatches(adds, 10).Error; err != nil {
		return err
	}
	return nil
}
func (r *addressRepo) GetAddressByUser(ctx context.Context, userID string) ([]entities.Address, error) {
	var addresses []entities.Address
	result := r.db.WithContext(ctx).Where("user_id=?", userID).Find(&addresses)
	if result.Error != nil {
		return nil, result.Error
	}
	return addresses, nil
}

func (r *addressRepo) ClearAllAddressesDataFromDB() error {
	if err := r.db.Exec("DELETE FROM addresses").Error; err != nil {
		return fmt.Errorf("failed to clear addresses table: %w", err)
	}
	return nil
}
