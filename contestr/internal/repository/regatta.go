package repository

import (
	"contestr/internal/configs"
	"contestr/pkg/regatta"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourRepository interface {
	Create(ctx context.Context, tour *regatta.Tour) (primitive.ObjectID, error)
	FindAll(ctx context.Context) ([]*regatta.Tour, error)
}

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

func (r *MongoTourRepository) FindByContestID(ctx context.Context, contestID int) (*regatta.Tour, error) {
	var tour regatta.Tour
	err := r.collection.FindOne(ctx, bson.M{"contest_id": contestID}).Decode(&tour)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &tour, nil
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

//	Tour: Tour{
//		Name:      tc.Name,
//		Index:     tc.Index,
//		StartTime: util.ParseTimeOrPanic(tc.StartTime),
//		Duration:  time.Duration(tc.Duration) * time.Minute,
//		Groups:    ConvertGroups(tc.Groups),
//		GroupSize: func() int {
//			if len(tc.Groups) == 0 {
//				return 0
//			}
//			return len(tc.Groups[0])
//		}(),
//		Problems:  tc.Problems,
//		ContestID: tc.ContestID,
//	},
//func ConvertGroups(groups [][]int) map[Participant]Group {
//	result := make(map[Participant]Group)
//
//	for _, group := range groups {
//		for _, participantID := range group {
//			result[participantID] = group
//		}
//	}
//
//	return result
//}
