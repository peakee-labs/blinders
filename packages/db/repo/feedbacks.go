package repo

import (
	"context"
	"time"

	"blinders/packages/db/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbacksRepo struct {
	Col *mongo.Collection
}

func NewFeedbackRepo(col *mongo.Collection) *FeedbacksRepo {
	return &FeedbacksRepo{Col: col}
}

func (r *FeedbacksRepo) InsertNewFeedback(f models.Feedback) (*models.Feedback, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	f.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := r.Col.InsertOne(ctx, f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
