package repository

import (
	"context"

	"contestr/internal/configs"
	"contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultTourTimetableCollection = "tour_timetables"

type TourTimetableRepository struct {
	collection *mongo.Collection
}

func NewTourTimetableRepository(
	cfg *configs.Config,
	client *mongo.Client,
) *TourTimetableRepository {
	collection := client.
		Database(cfg.MongoDB.Database).
		Collection(tourTimetableCollectionName(cfg))

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "contest_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, _ = collection.Indexes().CreateOne(context.Background(), indexModel)

	return &TourTimetableRepository{
		collection: collection,
	}
}

func tourTimetableCollectionName(cfg *configs.Config) string {
	if cfg.MongoDB.TourTimetableCollection != "" {
		return cfg.MongoDB.TourTimetableCollection
	}
	return defaultTourTimetableCollection
}

func (r *TourTimetableRepository) Create(ctx context.Context, timetable *regatta.ToursTimetable) error {
	_, err := r.collection.InsertOne(ctx, timetable)
	return err
}

func (r *TourTimetableRepository) GetByContestID(ctx context.Context, contestID int) (*regatta.ToursTimetable, error) {
	var timetable regatta.ToursTimetable
	err := r.collection.FindOne(ctx, bson.M{"contest_id": contestID}).Decode(&timetable)
	if err != nil {
		return nil, err
	}
	return &timetable, nil
}

func (r *TourTimetableRepository) Update(ctx context.Context, timetable *regatta.ToursTimetable) error {
	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"contest_id": timetable.ContestId},
		timetable,
	)
	return err
}

func (r *TourTimetableRepository) DeleteByContestID(ctx context.Context, contestID int) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"contest_id": contestID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
