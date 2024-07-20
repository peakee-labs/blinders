package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const FriendRequestsCollection = "friend-requests"

type FriendRequestsRepo struct {
	*mongo.Collection
}

func NewFriendRequestsRepo(db *mongo.Database) *FriendRequestsRepo {
	// TODO: need to add index for `to` and `status` for optimization,
	// mostly querying pending requests to someone
	return &FriendRequestsRepo{db.Collection(FriendRequestsCollection)}
}

func (r *FriendRequestsRepo) InsertNewRawFriendRequest(
	request FriendRequest,
) (*FriendRequest, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second)
	defer cancel()

	request.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	request.CreatedAt = now
	request.UpdatedAt = now

	upsert := true
	result, err := r.UpdateOne(ctx,
		bson.M{"from": request.From, "to": request.To},
		bson.M{"$setOnInsert": request},
		&options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		log.Println("can not insert friend request:", err)
		return nil, fmt.Errorf("something went wrong")
	} else if result.UpsertedCount == 0 {
		return nil, fmt.Errorf("request already existed")
	}

	return &request, err
}

func (r *FriendRequestsRepo) GetFriendRequestByFrom(
	from primitive.ObjectID,
	status FriendRequestStatus,
) ([]FriendRequest, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second)
	defer cancel()

	var filter bson.M
	switch status {
	case FriendStatusPending, FriendStatusAccepted, FriendStatusDenied:
		filter = bson.M{"from": from, "status": status}
	default:
		filter = bson.M{"from": from}
	}

	cursor, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var requests []FriendRequest
	err = cursor.All(ctx, &requests)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *FriendRequestsRepo) GetFriendRequestByTo(
	to primitive.ObjectID,
	status FriendRequestStatus,
) ([]FriendRequest, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second)
	defer cancel()

	var filter bson.M
	switch status {
	case FriendStatusPending, FriendStatusAccepted, FriendStatusDenied:
		filter = bson.M{"to": to, "status": status}
	default:
		filter = bson.M{"to": to}
	}

	cursor, err := r.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var requests []FriendRequest
	err = cursor.All(ctx, &requests)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *FriendRequestsRepo) GetFriendRequestByID(
	id primitive.ObjectID,
) (*FriendRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var request FriendRequest
	err := r.FindOne(ctx, bson.M{"_id": id}).Decode(&request)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("not found this friend request")
	} else if err != nil {
		log.Println("can not get friend request:", err)
		return nil, fmt.Errorf("something went wrong when get friend request")
	}

	return &request, nil
}

func (r *FriendRequestsRepo) UpdateFriendRequestStatusByID(
	id primitive.ObjectID,
	userID primitive.ObjectID,
	status FriendRequestStatus,
) (*FriendRequest, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(), time.Second)
	defer cancel()

	result, err := r.UpdateOne(
		ctx,
		bson.M{"_id": id, "to": userID, "status": FriendStatusPending},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		log.Println("can not update friend request:", err)
		return nil, fmt.Errorf("can not update friend request")
	}
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("not found this friend request")
	}
	var request FriendRequest
	err = r.FindOne(ctx, bson.M{"_id": id}).Decode(&request)
	if err != nil {
		log.Println("can not get friend request:", err)
		return nil, fmt.Errorf("something went wrong")
	}

	return &request, nil
}
