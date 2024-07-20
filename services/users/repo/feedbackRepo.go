package repo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackRepo struct {
	*mongo.Collection
}

func NewFeedbackRepo(db *mongo.Database) *FeedbackRepo {
	return &FeedbackRepo{db.Collection(FeedbackCollection)}
}

func (r *FeedbackRepo) InsertNewFeedback(f Feedback) (*Feedback, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	f.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	_, err := r.InsertOne(ctx, f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}
