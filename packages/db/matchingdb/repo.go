package matchingdb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MatchingCollection = "matching"

type MatchingRepo struct {
	*mongo.Collection
}

func NewMatchingRepo(db *mongo.Database) *MatchingRepo {
	col := db.Collection(MatchingCollection)
	ctx, cal := context.WithTimeout(context.Background(), time.Second*5)
	defer cal()

	if _, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"userId": 1},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		log.Println("can not create index for userId:", err)
		return nil
	}

	return &MatchingRepo{
		col,
	}
}

func (r *MatchingRepo) InsertNewRawMatchInfo(doc MatchInfo) (MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	_, err := r.InsertOne(ctx, doc)

	return doc, err
}

func (r *MatchingRepo) GetMatchInfoByUserID(userID primitive.ObjectID) (MatchInfo, error) {
	ctx, cal := context.WithTimeout(context.Background(), 5*time.Second)
	defer cal()

	var doc MatchInfo
	res := r.FindOne(ctx, bson.M{"userId": userID})
	if err := res.Err(); err != nil {
		return MatchInfo{}, err
	}
	if err := res.Decode(&doc); err != nil {
		return MatchInfo{}, err
	}
	return doc, nil
}

// GetUsersByLanguage returns `limit` ID of users that speak one language of `learnings` and are currently learning `native` or are currently learning same language as user.
func (r *MatchingRepo) GetUsersByLanguage(
	userID primitive.ObjectID,
	limit uint32,
) ([]string, error) {
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
					"native": bson.M{
						"$in": user.Learnings,
					}, // Users must speak at least one language of `learnings`.
					"learnings": bson.M{
						"$in": []string{user.Native},
					}, // Users should be learning their `native`.
				},
				{
					"learnings": bson.M{
						"$in": user.Learnings,
					}, // Users who learn the same language as the current user.
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

	cur, err := r.Aggregate(ctx, stages)
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

func (r *MatchingRepo) DropMatchInfoByUserID(userID primitive.ObjectID) (MatchInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	filter := bson.M{"userId": userID}
	res := r.FindOneAndDelete(ctx, filter)
	if err := res.Err(); err != nil {
		return MatchInfo{}, err
	}
	var deletedUser MatchInfo
	if err := res.Decode(&deletedUser); err != nil {
		return MatchInfo{}, err
	}
	return deletedUser, nil
}
