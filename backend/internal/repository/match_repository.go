package repository

import (
	"context"

	"poke/backend/internal/model"
)

// MatchRepository 定义对局仓储接口。
type MatchRepository interface {
	Create(ctx context.Context, match *model.Match) error
	Update(ctx context.Context, match *model.Match) error
	FindByID(ctx context.Context, id string) (*model.Match, error)
	FindByRoomCode(ctx context.Context, roomCode string) (*model.Match, error)
	ListByPlayer(ctx context.Context, userID string, limit, offset int64) ([]*model.Match, error)
}

