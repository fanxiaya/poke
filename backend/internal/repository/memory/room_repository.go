package memory

import (
	"context"
	"sync"

	"poke/backend/internal/model"
	"poke/backend/internal/repository"
)

// RoomRepository 为房间仓储提供内存实现。
type RoomRepository struct {
	mu     sync.RWMutex
	byID   map[string]*model.Room
	byCode map[string]string
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{
		byID:   make(map[string]*model.Room),
		byCode: make(map[string]string),
	}
}

func (r *RoomRepository) Create(_ context.Context, room *model.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if room == nil {
		return nil
	}
	r.byID[room.ID] = cloneRoom(room)
	r.byCode[room.RoomCode] = room.ID
	return nil
}

func (r *RoomRepository) Update(_ context.Context, room *model.Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if room == nil {
		return nil
	}
	if _, ok := r.byID[room.ID]; !ok {
		return repository.ErrNotFound
	}
	r.byID[room.ID] = cloneRoom(room)
	r.byCode[room.RoomCode] = room.ID
	return nil
}

func (r *RoomRepository) FindByID(_ context.Context, id string) (*model.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneRoom(room), nil
}

func (r *RoomRepository) FindByCode(_ context.Context, roomCode string) (*model.Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byCode[roomCode]
	if !ok {
		return nil, repository.ErrNotFound
	}
	room, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneRoom(room), nil
}

