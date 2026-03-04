package mongo

import (
	"context"

	"poke/backend/internal/model"
	"poke/backend/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// RoomRepository 为房间仓储提供 MongoDB 实现。
type RoomRepository struct {
	collection *mongo.Collection
}

func NewRoomRepository(db *mongo.Database) *RoomRepository {
	return &RoomRepository{
		collection: db.Collection(roomsCollection),
	}
}

func (r *RoomRepository) Create(ctx context.Context, room *model.Room) error {
	if room == nil {
		return nil
	}
	if room.ID == "" {
		room.ID = bson.NewObjectID().Hex()
	}
	_, err := r.collection.InsertOne(ctx, room)
	return err
}

func (r *RoomRepository) Update(ctx context.Context, room *model.Room) error {
	if room == nil {
		return nil
	}
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": room.ID}, room)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *RoomRepository) FindByID(ctx context.Context, id string) (*model.Room, error) {
	var room model.Room
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepository) FindByCode(ctx context.Context, roomCode string) (*model.Room, error) {
	var room model.Room
	err := r.collection.FindOne(ctx, bson.M{"roomCode": roomCode}).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &room, nil
}
