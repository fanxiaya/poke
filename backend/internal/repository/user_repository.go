package repository

import (
	"context"

	"poke/backend/internal/model"
)

// UserRepository 定义用户仓储接口。
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByOpenID(ctx context.Context, openID string) (*model.User, error)
}

