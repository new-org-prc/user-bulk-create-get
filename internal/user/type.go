package user

import (
	"context"
	"sika/pkg/storage/entities"
)


type Repo interface{
	CreateUser(ctx context.Context, user *entities.User)error
	CreateBatchUsers(ctx context.Context, users []entities.User)error
	GetUserByID(ctx context.Context, id string)(*entities.User, error)
	ClearAllUsersDataFromDB()error
}