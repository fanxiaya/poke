package mongo

import "go.mongodb.org/mongo-driver/v2/mongo"

const (
	usersCollection   = "users"
	roomsCollection   = "rooms"
	matchesCollection = "matches"
)

// RepositorySet 聚合 Mongo 仓储实现，方便业务层注入。
type RepositorySet struct {
	User  *UserRepository
	Room  *RoomRepository
	Match *MatchRepository
}

func NewRepositorySet(db *mongo.Database) *RepositorySet {
	return &RepositorySet{
		User:  NewUserRepository(db),
		Room:  NewRoomRepository(db),
		Match: NewMatchRepository(db),
	}
}

