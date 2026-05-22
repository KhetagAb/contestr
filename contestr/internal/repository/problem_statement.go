package repository

import (
	"context"
	"time"

	"contestr/internal/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProblemStatement struct {
	ContestID   int       `bson:"contest_id"`
	ProblemCode string    `bson:"problem_code"`
	ObjectKey   string    `bson:"object_key"`
	UploadedAt  time.Time `bson:"uploaded_at"`
	SizeBytes   int64     `bson:"size_bytes,omitempty"`
}

type ProblemStatementRepository interface {
	Upsert(ctx context.Context, doc ProblemStatement) error
	Get(ctx context.Context, contestID int, problemCode string) (*ProblemStatement, error)
	ListByContest(ctx context.Context, contestID int) ([]ProblemStatement, error)
	Delete(ctx context.Context, contestID int, problemCode string) error
	DeleteByContest(ctx context.Context, contestID int) error
}

type MongoProblemStatementRepository struct {
	collection *mongo.Collection
}

func NewMongoProblemStatementRepository(cfg *configs.Config, client *mongo.Client) *MongoProblemStatementRepository {
	collection := client.Database(cfg.MongoDB.Database).Collection("problem_statements")
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "contest_id", Value: 1},
			{Key: "problem_code", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}
	_, _ = collection.Indexes().CreateOne(context.Background(), indexModel)
	return &MongoProblemStatementRepository{collection: collection}
}

func (r *MongoProblemStatementRepository) Upsert(ctx context.Context, doc ProblemStatement) error {
	filter := bson.M{"contest_id": doc.ContestID, "problem_code": doc.ProblemCode}
	update := bson.M{"$set": doc}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *MongoProblemStatementRepository) Get(ctx context.Context, contestID int, problemCode string) (*ProblemStatement, error) {
	var doc ProblemStatement
	err := r.collection.FindOne(ctx, bson.M{
		"contest_id":   contestID,
		"problem_code": problemCode,
	}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &doc, nil
}

func (r *MongoProblemStatementRepository) ListByContest(ctx context.Context, contestID int) ([]ProblemStatement, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"contest_id": contestID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []ProblemStatement
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	if docs == nil {
		return []ProblemStatement{}, nil
	}
	return docs, nil
}

func (r *MongoProblemStatementRepository) Delete(ctx context.Context, contestID int, problemCode string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{
		"contest_id":   contestID,
		"problem_code": problemCode,
	})
	return err
}

func (r *MongoProblemStatementRepository) DeleteByContest(ctx context.Context, contestID int) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"contest_id": contestID})
	return err
}
