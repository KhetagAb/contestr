package repository

import (
	"context"
	"errors"
	"time"

	"contestr/internal/configs"
	"contestr/pkg/regatta"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RegisteredContest struct {
	ContestID       int                     `bson:"contest_id"`
	System          string                  `bson:"system"`
	Name            string                  `bson:"name"`
	CreatedAt       time.Time               `bson:"created_at"`
	ScoringSettings regatta.ScoringSettings `bson:"scoring_settings"`
	TourSettings    regatta.TourSettings    `bson:"tour_settings"`
}

type RegisteredContestRepository interface {
	List(ctx context.Context) ([]RegisteredContest, error)
	GetByContestID(ctx context.Context, contestID int) (*RegisteredContest, error)
	Create(ctx context.Context, contest RegisteredContest) error
	UpdateContestSettings(ctx context.Context, contestID int, scoring regatta.ScoringSettings, tour regatta.TourSettings) error
	Delete(ctx context.Context, contestID int) error
}

type MongoRegisteredContestRepository struct {
	collection *mongo.Collection
}

func NewMongoRegisteredContestRepository(
	cfg *configs.Config,
	client *mongo.Client,
) *MongoRegisteredContestRepository {
	collection := client.
		Database(cfg.MongoDB.Database).
		Collection("registered_contests")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "contest_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, _ = collection.Indexes().CreateOne(context.Background(), indexModel)

	return &MongoRegisteredContestRepository{collection: collection}
}

func (r *MongoRegisteredContestRepository) List(ctx context.Context) ([]RegisteredContest, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "contest_id", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var contests []RegisteredContest
	if err := cursor.All(ctx, &contests); err != nil {
		return nil, err
	}
	if contests == nil {
		return []RegisteredContest{}, nil
	}
	for i := range contests {
		contests[i].ScoringSettings = regatta.NormalizeScoringSettings(contests[i].ScoringSettings)
		contests[i].TourSettings = regatta.NormalizeTourSettings(contests[i].TourSettings)
	}
	return contests, nil
}

func (r *MongoRegisteredContestRepository) GetByContestID(ctx context.Context, contestID int) (*RegisteredContest, error) {
	var contest RegisteredContest
	err := r.collection.FindOne(ctx, bson.M{"contest_id": contestID}).Decode(&contest)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	contest.ScoringSettings = regatta.NormalizeScoringSettings(contest.ScoringSettings)
	contest.TourSettings = regatta.NormalizeTourSettings(contest.TourSettings)
	return &contest, nil
}

func (r *MongoRegisteredContestRepository) Create(ctx context.Context, contest RegisteredContest) error {
	contest.ScoringSettings = regatta.NormalizeScoringSettings(contest.ScoringSettings)
	contest.TourSettings = regatta.NormalizeTourSettings(contest.TourSettings)
	_, err := r.collection.InsertOne(ctx, contest)
	if mongo.IsDuplicateKeyError(err) {
		return ErrContestAlreadyRegistered
	}
	return err
}

func (r *MongoRegisteredContestRepository) UpdateContestSettings(ctx context.Context, contestID int, scoring regatta.ScoringSettings, tour regatta.TourSettings) error {
	scoring = regatta.NormalizeScoringSettings(scoring)
	tour = regatta.NormalizeTourSettings(tour)
	res, err := r.collection.UpdateOne(
		ctx,
		bson.M{"contest_id": contestID},
		bson.M{"$set": bson.M{
			"scoring_settings": scoring,
			"tour_settings":    tour,
		}},
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrContestNotRegistered
	}
	return nil
}

func (r *MongoRegisteredContestRepository) Delete(ctx context.Context, contestID int) error {
	res, err := r.collection.DeleteOne(ctx, bson.M{"contest_id": contestID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrContestNotRegistered
	}
	return nil
}
