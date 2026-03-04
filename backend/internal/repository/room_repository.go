package repository

import (
	"context"

	"poke/backend/internal/model"
)

// RoomRepository 定义房间仓储接口。
type RoomRepository interface {
	Create(ctx context.Context, room *model.Room) error
	Update(ctx context.Context, room *model.Room) error
	FindByID(ctx context.Context, id string) (*model.Room, error)
	FindByCode(ctx context.Context, roomCode string) (*model.Room, error)
}

