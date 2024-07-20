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

const UsersCollection = "users"

type UsersRepo struct {
	*mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	col := db.Collection(UsersCollection)
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	_, err := col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"firebaseUID": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Println("can not create index for firebaseUID:", err)
		return nil
	}

	return &UsersRepo{col}
}

func (r *UsersRepo) InsertNewUser(u User) (User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	_, err := r.InsertOne(ctx, u)

	return u, err
}

// this function creates new ID and time and insert the document to database
func (r *UsersRepo) InsertNewRawUser(u User) (User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	u.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	u.CreatedAt = now
	u.UpdatedAt = now

	_, err := r.InsertOne(ctx, u)

	return u, err
}

func (r *UsersRepo) GetUserByID(id primitive.ObjectID) (User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var user User
	err := r.FindOne(ctx, bson.M{"_id": id}).Decode(&user)

	return user, err
}

func (r *UsersRepo) GetUserByFirebaseUID(firebaseUID string) (User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var user User
	err := r.FindOne(ctx, bson.M{"firebaseUID": firebaseUID}).Decode(&user)

	return user, err
}

func (r *UsersRepo) GetUserByEmail(email string) (User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var user User
	err := r.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	return user, err
}

func (r *UsersRepo) DeleteUserByID(userID primitive.ObjectID) (User, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	usr := User{}
	err := r.FindOneAndDelete(ctx, bson.M{"_id": userID}).Decode(&usr)
	return usr, err
}

func (r *UsersRepo) AddFriend(user1ID primitive.ObjectID, user2ID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	result, err := r.BulkWrite(
		ctx,
		[]mongo.WriteModel{
			mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": user1ID}).
				SetUpdate(bson.M{"$addToSet": bson.M{"friends": user2ID}}),
			mongo.NewUpdateOneModel().
				SetFilter(bson.M{"_id": user2ID}).
				SetUpdate(bson.M{"$addToSet": bson.M{"friends": user1ID}}),
		},
	)
	if err != nil {
		log.Println("can not add friend:", err)
		return fmt.Errorf("something went wrong")
	} else if result.ModifiedCount != 2 {
		log.Println("wrong updated count when add friend")
		return fmt.Errorf("update friend failed, wrong updated count")
	}

	return nil
}
