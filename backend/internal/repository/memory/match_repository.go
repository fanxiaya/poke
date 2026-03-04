package memory

import (
	"context"
	"sort"
	"sync"

	"poke/backend/internal/model"
	"poke/backend/internal/repository"
)

// MatchRepository 为对局仓储提供内存实现。
type MatchRepository struct {
	mu         sync.RWMutex
	byID       map[string]*model.Match
	byRoomCode map[string]string
}

func NewMatchRepository() *MatchRepository {
	return &MatchRepository{
		byID:       make(map[string]*model.Match),
		byRoomCode: make(map[string]string),
	}
}

func (r *MatchRepository) Create(_ context.Context, match *model.Match) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if match == nil {
		return nil
	}
	r.byID[match.ID] = cloneMatch(match)
	r.byRoomCode[match.RoomCode] = match.ID
	return nil
}

func (r *MatchRepository) Update(_ context.Context, match *model.Match) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if match == nil {
		return nil
	}
	if _, ok := r.byID[match.ID]; !ok {
		return repository.ErrNotFound
	}
	r.byID[match.ID] = cloneMatch(match)
	r.byRoomCode[match.RoomCode] = match.ID
	return nil
}

func (r *MatchRepository) FindByID(_ context.Context, id string) (*model.Match, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	match, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneMatch(match), nil
}

func (r *MatchRepository) FindByRoomCode(_ context.Context, roomCode string) (*model.Match, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byRoomCode[roomCode]
	if !ok {
		return nil, repository.ErrNotFound
	}
	match, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneMatch(match), nil
}

func (r *MatchRepository) ListByPlayer(_ context.Context, userID string, limit, offset int64) ([]*model.Match, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*model.Match, 0)
	for _, m := range r.byID {
		if containsMatchPlayer(m.Players, userID) {
			filtered = append(filtered, cloneMatch(m))
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = int64(len(filtered))
	}

	if offset >= int64(len(filtered)) {
		return []*model.Match{}, nil
	}
	end := offset + limit
	if end > int64(len(filtered)) {
		end = int64(len(filtered))
	}
	return filtered[offset:end], nil
}

func containsMatchPlayer(players []model.MatchPlayer, userID string) bool {
	for _, player := range players {
		if player.UserID == userID {
			return true
		}
	}
	return false
}

