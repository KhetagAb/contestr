package repository

import (
	"contestr/internal/configs"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CodeforcesHandleMapping struct {
	ContestID string `bson:"contest_id"`
	Handle    string `bson:"handle"`
	Name      string `bson:"name"`
}

type CodeforcesHandleRepository interface {
	GetAllByContestID(ctx context.Context, contestID int) (map[string]string, error)
	ListByContestID(ctx context.Context, contestID int) ([]CodeforcesHandleMapping, error)
	UpsertMany(ctx context.Context, contestID int, mappings []CodeforcesHandleMapping) error
	DeleteOne(ctx context.Context, contestID int, handle string) error
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
	_, _ = collection.Indexes().CreateOne(context.Background(), indexModel)

	return &MongoCodeforcesHandleRepository{
		collection: collection,
	}
}

func contestIDString(contestID int) string {
	return fmt.Sprintf("%d", contestID)
}

func (r *MongoCodeforcesHandleRepository) GetAllByContestID(ctx context.Context, contestID int) (map[string]string, error) {
	mappings, err := r.ListByContestID(ctx, contestID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(mappings))
	for _, mapping := range mappings {
		result[mapping.Handle] = mapping.Name
	}
	return result, nil
}

func (r *MongoCodeforcesHandleRepository) ListByContestID(ctx context.Context, contestID int) ([]CodeforcesHandleMapping, error) {
	filter := bson.M{"contest_id": contestIDString(contestID)}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "handle", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mappings []CodeforcesHandleMapping
	if err = cursor.All(ctx, &mappings); err != nil {
		return nil, err
	}
	if mappings == nil {
		return []CodeforcesHandleMapping{}, nil
	}
	return mappings, nil
}

func (r *MongoCodeforcesHandleRepository) UpsertMany(ctx context.Context, contestID int, mappings []CodeforcesHandleMapping) error {
	contestIDStr := contestIDString(contestID)
	for _, m := range mappings {
		if m.Handle == "" {
			continue
		}
		name := m.Name
		if name == "" {
			name = m.Handle
		}
		filter := bson.M{"contest_id": contestIDStr, "handle": m.Handle}
		update := bson.M{"$set": bson.M{
			"contest_id": contestIDStr,
			"handle":     m.Handle,
			"name":       name,
		}}
		opts := options.Update().SetUpsert(true)
		if _, err := r.collection.UpdateOne(ctx, filter, update, opts); err != nil {
			return err
		}
	}
	return nil
}

func (r *MongoCodeforcesHandleRepository) DeleteOne(ctx context.Context, contestID int, handle string) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{
		"contest_id": contestIDString(contestID),
		"handle":     handle,
	})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrHandleNotFound
	}
	return nil
}
