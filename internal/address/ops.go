package address

import (
	"context"
	"sika/pkg/storage/entities"
)

type Ops struct {
	repo Repo
}

func NewOps(repo Repo) *Ops {
	return &Ops{repo}
}

func (o *Ops) CreateAddress(ctx context.Context, a *entities.Address) error {
	return o.repo.CreateAddress(ctx, a)
}

func (o *Ops) CreateBatchAddress(ctx context.Context, addrs []entities.Address) error {
	return o.repo.CreateBatchAddresses(ctx, addrs)
}

func (o *Ops) GetAddressByUserID(ctx context.Context, uid string) ([]entities.Address, error) {
	return o.repo.GetAddressByUser(ctx, uid)
}

func (o *Ops) ClearAllAddressesDataFromDB() error {
	return o.repo.ClearAllAddressesDataFromDB()
}
