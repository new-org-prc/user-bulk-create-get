package user

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


func (o *Ops) CreateUser(ctx context.Context, user *entities.User) error {
	return o.repo.CreateUser(ctx, user)
}

func (o *Ops) CreateBatchUser(ctx context.Context, users []entities.User) error {
	return o.repo.CreateBatchUsers(ctx, users)
}

func (o *Ops) GetUserByID(ctx context.Context, uid string) (*entities.User, error) {
	return o.repo.GetUserByID(ctx, uid)
}

func (o *Ops) ClearAllUsersDataFromDB() error {
	return o.repo.ClearAllUsersDataFromDB()
}
