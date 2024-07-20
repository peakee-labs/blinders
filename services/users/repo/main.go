package repo

import "go.mongodb.org/mongo-driver/mongo"

var (
	UsersCollection          = "users"
	FriendRequestsCollection = "friend-requests"
	FeedbackCollection       = "feedback"
)

type UsersDB struct {
	mongo.Database
	UsersRepo          *UsersRepo
	FriendRequestsRepo *FriendRequestsRepo
	FeedbackRepo       *FeedbackRepo
}

func NewUsersDB(db *mongo.Database) *UsersDB {
	return &UsersDB{
		Database:           *db,
		UsersRepo:          NewUsersRepo(db),
		FriendRequestsRepo: NewFriendRequestsRepo(db),
		FeedbackRepo:       NewFeedbackRepo(db),
	}
}
