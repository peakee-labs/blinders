package repo

import "go.mongodb.org/mongo-driver/bson/primitive"

type ConversationType string

const (
	IndividualConversation ConversationType = "individual"
	GroupConversation      ConversationType = "group"
)

type Conversation struct {
	ID        primitive.ObjectID    `bson:"_id"                json:"id"`
	Type      ConversationType      `bson:"type"               json:"type"`
	Members   []Member              `bson:"members"            json:"members"`
	CreatedBy primitive.ObjectID    `bson:"createdBy"          json:"createdBy"`
	CreatedAt primitive.DateTime    `bson:"createdAt"          json:"createdAt"`
	UpdatedAt primitive.DateTime    `bson:"updatedAt"          json:"updatedAt"`
	Metadata  *ConversationMetadata `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

type ConversationMetadata struct {
	Name  string `bson:"name,omitempty"  json:"name,omitempty"`
	Image string `bson:"image,omitempty" json:"image,omitempty"`
}

type Member struct {
	UserID                primitive.ObjectID  `bson:"userId"                          json:"userId"`
	Nickname              string              `bson:"nickname,omitempty"              json:"nickname,omitempty"`
	LatestViewedMessageID *primitive.ObjectID `bson:"latestViewedMessageId,omitempty" json:"latestViewedMessageId,omitempty"`
	CreatedAt             primitive.DateTime  `bson:"createdAt"                       json:"createdAt"`
	UpdatedAt             primitive.DateTime  `bson:"updatedAt"                       json:"updatedAt"`
	JoinedAt              primitive.DateTime  `bson:"joinedAt"                        json:"joinedAt"`
}

type MessageStatus string

const (
	DeliveredStatus MessageStatus = "delivered"
	ReceivedStatus  MessageStatus = "received"
	SeenStatus      MessageStatus = "seen"
)

type Message struct {
	ID             primitive.ObjectID  `bson:"_id"               json:"id"`
	SenderID       primitive.ObjectID  `bson:"senderId"          json:"senderId"`
	ConversationID primitive.ObjectID  `bson:"conversationId"    json:"conversationId"`
	ReplyTo        *primitive.ObjectID `bson:"replyTo,omitempty" json:"replyTo,omitempty"`
	Content        string              `bson:"content"           json:"content"`
	Status         MessageStatus       `bson:"status"            json:"status"`
	CreatedAt      primitive.DateTime  `bson:"createdAt"         json:"createdAt"`
	UpdatedAt      primitive.DateTime  `bson:"updatedAt"         json:"updatedAt"`
	Emotions       []MessageEmotion    `bson:"emotions"          json:"emotions"`
}

type MessageEmotion struct {
	SenderID  primitive.ObjectID `bson:"senderId"  json:"senderId"`
	Content   string             `bson:"content"   json:"content"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}
