package repository

import (
	"contestr/internal/configs"
	"contestr/pkg/regatta"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContestRepository interface {
	Upsert(ctx context.Context, data *regatta.Contest) error
	GetByContestID(ctx context.Context, contestID int) (*regatta.Contest, error)
	GetParticipants(ctx context.Context, contestID int) (map[string]string, error)
	GetSubmissions(ctx context.Context, contestID int) ([]regatta.ContestSubmission, error)
	DeleteByContestID(ctx context.Context, contestID int) error
}

type MongoContestRepository struct {
	collection *mongo.Collection
}

func NewMongoContestRepository(
	cfg *configs.Config,
	client *mongo.Client,
) *MongoContestRepository {
	collection := client.
		Database(cfg.MongoDB.Database).
		Collection("contests")
	return &MongoContestRepository{
		collection: collection,
	}
}

func (r *MongoContestRepository) Upsert(ctx context.Context, data *regatta.Contest) error {
	filter := bson.M{"contest_id": data.ContestID}
	update := bson.M{"$set": data}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *MongoContestRepository) GetByContestID(ctx context.Context, contestID int) (*regatta.Contest, error) {
	var contest regatta.Contest
	err := r.collection.FindOne(ctx, bson.M{"contest_id": contestID}).Decode(&contest)
	if err != nil {
		return nil, err
	}
	contest.ScoringSettings = regatta.NormalizeScoringSettings(contest.ScoringSettings)
	contest.TourSettings = regatta.NormalizeTourSettings(contest.TourSettings)
	return &contest, nil
}

func (r *MongoContestRepository) GetParticipants(ctx context.Context, contestID int) (map[string]string, error) {
	contest, err := r.GetByContestID(ctx, contestID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, participant := range contest.Participants {
		result[participant.ID] = participant.DisplayName
	}
	return result, nil
}

func (r *MongoContestRepository) GetSubmissions(ctx context.Context, contestID int) ([]regatta.ContestSubmission, error) {
	contest, err := r.GetByContestID(ctx, contestID)
	if err != nil {
		return nil, err
	}
	return contest.Submissions, nil
}

func (r *MongoContestRepository) DeleteByContestID(ctx context.Context, contestID int) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"contest_id": contestID})
	return err
}
