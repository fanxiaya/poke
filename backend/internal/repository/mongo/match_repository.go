package mongo

import (
	"context"

	"poke/backend/internal/model"
	"poke/backend/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MatchRepository 为对局仓储提供 MongoDB 实现。
type MatchRepository struct {
	collection *mongo.Collection
}

func NewMatchRepository(db *mongo.Database) *MatchRepository {
	return &MatchRepository{
		collection: db.Collection(matchesCollection),
	}
}

func (r *MatchRepository) Create(ctx context.Context, match *model.Match) error {
	if match == nil {
		return nil
	}
	if match.ID == "" {
		match.ID = bson.NewObjectID().Hex()
	}
	_, err := r.collection.InsertOne(ctx, match)
	return err
}

func (r *MatchRepository) Update(ctx context.Context, match *model.Match) error {
	if match == nil {
		return nil
	}
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": match.ID}, match)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *MatchRepository) FindByID(ctx context.Context, id string) (*model.Match, error) {
	var match model.Match
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&match)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) FindByRoomCode(ctx context.Context, roomCode string) (*model.Match, error) {
	var match model.Match
	err := r.collection.FindOne(ctx, bson.M{"roomCode": roomCode}).Decode(&match)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &match, nil
}

func (r *MatchRepository) ListByPlayer(ctx context.Context, userID string, limit, offset int64) ([]*model.Match, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	findOptions := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetLimit(limit).
		SetSkip(offset)

	cursor, err := r.collection.Find(ctx, bson.M{"players.userId": userID}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make([]*model.Match, 0)
	for cursor.Next(ctx) {
		var match model.Match
		if decodeErr := cursor.Decode(&match); decodeErr != nil {
			return nil, decodeErr
		}
		result = append(result, &match)
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
