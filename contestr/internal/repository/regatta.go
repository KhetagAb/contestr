package repository

import (
	"contestr/internal/configs"
	"contestr/pkg/regatta"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTourRepository struct {
	collection *mongo.Collection
}

func NewMongoTourRepository(
	cfg *configs.Config,
	client *mongo.Client,
) *MongoTourRepository {
	collection := client.
		Database(cfg.MongoDB.Database).
		Collection(cfg.MongoDB.TourCollection)
	return &MongoTourRepository{
		collection: collection,
	}
}

func (r *MongoTourRepository) Create(ctx context.Context, tour *regatta.Tour) (primitive.ObjectID, error) {
	result, err := r.collection.InsertOne(ctx, tour)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *MongoTourRepository) UpdateDuration(ctx context.Context, contestID int, sequence int, durationSeconds int) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"contest_id": contestID, "sequence": sequence},
		bson.M{"$set": bson.M{"duration_in_seconds": durationSeconds}},
	)
	return err
}

func (r *MongoTourRepository) FindByContestID(ctx context.Context, contestID int) ([]regatta.Tour, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"contest_id": contestID},
		options.Find().SetSort(bson.D{{Key: "sequence", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tours []regatta.Tour
	if err = cursor.All(ctx, &tours); err != nil {
		return nil, err
	}
	return tours, nil
}

func (r *MongoTourRepository) FindAll(ctx context.Context) ([]regatta.Tour, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tours []regatta.Tour
	if err = cursor.All(ctx, &tours); err != nil {
		return nil, err
	}
	return tours, nil
}
