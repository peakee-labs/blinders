package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID   `bson:"_id"         json:"id"`
	Name        string               `bson:"name"        json:"name"`
	Email       string               `bson:"email"       json:"email"`
	FirebaseUID string               `bson:"firebaseUID" json:"firebaseUID"`
	ImageURL    string               `bson:"imageURL"    json:"imageURL"`
	FriendIDs   []primitive.ObjectID `bson:"friends"     json:"friends"`
	CreatedAt   primitive.DateTime   `bson:"createdAt"   json:"createdAt"`
	UpdatedAt   primitive.DateTime   `bson:"updatedAt"   json:"updatedAt"`
	// Conversations []EmbeddedConversation `bson:"conversations" json:"conversations"`
}

// not use embedded conversation now
// we could optimize conversation query by this later
// also we can add more fields to embedded conversation like individual settings

// type EmbeddedConversation struct {
// 	ID             primitive.ObjectID `bson:"_id"            json:"id"`
// 	ConversationID primitive.ObjectID `bson:"conversationId" json:"conversationId"`
// 	CreatedAt      primitive.DateTime `bson:"createdAt"      json:"createdAt"`
// 	UpdatedAt      primitive.DateTime `bson:"updatedAt"      json:"updatedAt"`
// 	Settings       struct {
// 		Notification bool `bson:"notification"   json:"notification"`
// 	} `bson:"settings"       json:"settings"`
// }

type FriendRequestStatus string

const (
	FriendStatusPending  FriendRequestStatus = "pending"
	FriendStatusAccepted FriendRequestStatus = "accepted"
	FriendStatusDenied   FriendRequestStatus = "denied"
)

type FriendRequest struct {
	ID        primitive.ObjectID  `bson:"_id"       json:"id"`
	From      primitive.ObjectID  `bson:"from"      json:"from"`
	To        primitive.ObjectID  `bson:"to"        json:"to"`
	Status    FriendRequestStatus `bson:"status"    json:"status"`
	CreatedAt primitive.DateTime  `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime  `bson:"updatedAt" json:"updatedAt"`
}
