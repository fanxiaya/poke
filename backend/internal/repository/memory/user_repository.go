package memory

import (
	"context"
	"sync"

	"poke/backend/internal/model"
	"poke/backend/internal/repository"
)

// UserRepository 为用户仓储提供内存实现。
type UserRepository struct {
	mu       sync.RWMutex
	byID     map[string]*model.User
	byOpenID map[string]string
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		byID:     make(map[string]*model.User),
		byOpenID: make(map[string]string),
	}
}

func (r *UserRepository) Create(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user == nil {
		return nil
	}
	r.byID[user.ID] = cloneUser(user)
	r.byOpenID[user.OpenID] = user.ID
	return nil
}

func (r *UserRepository) Update(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user == nil {
		return nil
	}
	if _, ok := r.byID[user.ID]; !ok {
		return repository.ErrNotFound
	}
	r.byID[user.ID] = cloneUser(user)
	r.byOpenID[user.OpenID] = user.ID
	return nil
}

func (r *UserRepository) FindByID(_ context.Context, id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneUser(user), nil
}

func (r *UserRepository) FindByOpenID(_ context.Context, openID string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byOpenID[openID]
	if !ok {
		return nil, repository.ErrNotFound
	}
	user, ok := r.byID[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return cloneUser(user), nil
}

