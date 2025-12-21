package repository

import (
	"contestr/internal/configs"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type codeforcesHandleMapping struct {
	ContestID string `bson:"contest_id"`
	Handle    string `bson:"handle"`
	Name      string `bson:"name"`
}

type CodeforcesHandleRepository interface {
	GetAllByContestID(ctx context.Context, contestID int) (map[string]string, error)
}

type MongoCodeforcesHandleRepository struct {
	collection *mongo.Collection
}

func NewMongoCodeforcesHandleRepository(
	cfg *configs.Config,
	client *mongo.Client,
) *MongoCodeforcesHandleRepository {
	collection := client.
		Database(cfg.MongoDB.Database).
		Collection("codeforces_handles")
	
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "contest_id", Value: 1},
			{Key: "handle", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	collection.Indexes().CreateOne(context.Background(), indexModel)
	
	return &MongoCodeforcesHandleRepository{
		collection: collection,
	}
}

func (r *MongoCodeforcesHandleRepository) GetAllByContestID(ctx context.Context, contestID int) (map[string]string, error) {
	filter := bson.M{"contest_id": fmt.Sprintf("%d", contestID)}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mappings []codeforcesHandleMapping
	if err = cursor.All(ctx, &mappings); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, mapping := range mappings {
		result[mapping.Handle] = mapping.Name
	}

	return result, nil
}
