package address

import (
	"context"
	"sika/pkg/storage/entities"
)

type Repo interface {
	CreateAddress(ctx context.Context, a *entities.Address)error
	CreateBatchAddresses(ctx context.Context, adds []entities.Address)error
	GetAddressByUser(ctx context.Context, userID string)([]entities.Address, error)
	ClearAllAddressesDataFromDB()error
}
