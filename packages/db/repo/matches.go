package repo

import (
	"context"
	"log"
	"time"

	"blinders/packages/db/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MatchesRepo struct {
	Col *mongo.Collection
}

func NewMatchesRepo(col *mongo.Collection) *MatchesRepo {
	ctx, cal := context.WithTimeout(context.Background(), time.Second*5)
	defer cal()

	if _, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"userId": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		log.Println("can not create index for userId:", err)
		return nil
	}

	return &MatchesRepo{
		Col: col,
	}
}

func (r *MatchesRepo) InsertNewRawMatchInfo(doc models.MatchInfo) (models.MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	_, err := r.Col.InsertOne(ctx, doc)

	return doc, err
}

func (r *MatchesRepo) GetMatchInfoByUserID(userID primitive.ObjectID) (models.MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	var doc models.MatchInfo
	res := r.Col.FindOne(ctx, bson.M{"userId": userID})
	if err := res.Err(); err != nil {
		return models.MatchInfo{}, err
	}
	if err := res.Decode(&doc); err != nil {
		return models.MatchInfo{}, err
	}
	return doc, nil
}

// GetUsersByLanguage returns `limit` ID of users that speak one language of `learnings` and are currently learning `native` or are currently learning same language as user.
func (r *MatchesRepo) GetUsersByLanguage(userID primitive.ObjectID, limit uint32) ([]string, error) {
	user, err := r.GetMatchInfoByUserID(userID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	stages := []bson.M{
		{"$match": bson.M{
			"userId": bson.M{"$ne": user.UserID},
			"$or": []bson.M{
				{
					"native":    bson.M{"$in": user.Learnings},        // Users must speak at least one language of `learnings`.
					"learnings": bson.M{"$in": []string{user.Native}}, // Users should be learning their `native`.
				},
				{
					"learnings": bson.M{"$in": user.Learnings}, // Users who learn the same language as the current user.
				},
			},
		}},
		// at here we may sort users based on any ranking mark from the system.
		// currently, we random choose 1000 user.
		{
			"$sample": bson.M{"size": limit},
		},
		{"$project": bson.M{"_id": 0, "userId": 1}},
	}

	cur, err := r.Col.Aggregate(ctx, stages)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = cur.Close(ctx); err != nil {
			log.Panicf("repo: cannot close cursor, err: %v", err)
		}
	}()

	type ReturnType struct {
		UserID primitive.ObjectID `bson:"userId"`
	}

	var ids []string
	for cur.Next(ctx) {
		doc := new(ReturnType)
		if err := cur.Decode(doc); err != nil {
			return nil, err
		}
		ids = append(ids, doc.UserID.Hex())
	}
	return ids, nil
}

func (r *MatchesRepo) DropMatchInfoByUserID(userID primitive.ObjectID) (models.MatchInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"userId": userID}
	res := r.Col.FindOneAndDelete(ctx, filter)
	if err := res.Err(); err != nil {
		return models.MatchInfo{}, err
	}
	var deletedUser models.MatchInfo
	if err := res.Decode(&deletedUser); err != nil {
		return models.MatchInfo{}, err
	}
	return deletedUser, nil
}
