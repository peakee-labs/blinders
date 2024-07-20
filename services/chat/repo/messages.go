package repo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MessagesCollection = "messages"

type MessagesRepo struct {
	*mongo.Collection
}

func NewMessagesRepo(db *mongo.Database) *MessagesRepo {
	return &MessagesRepo{db.Collection(MessagesCollection)}
}

func (r MessagesRepo) ConstructNewMessage(
	senderID primitive.ObjectID,
	conversationID primitive.ObjectID,
	replyTo primitive.ObjectID,
	content string,
) Message {
	now := primitive.NewDateTimeFromTime(time.Now())
	replyToPointer := &replyTo
	if replyTo.IsZero() {
		replyToPointer = nil
	}
	return Message{
		ID:             primitive.NewObjectID(),
		Status:         "delivered",
		Emotions:       make([]MessageEmotion, 0),
		SenderID:       senderID,
		ConversationID: conversationID,
		ReplyTo:        replyToPointer,
		Content:        content,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func (r *MessagesRepo) GetMessageByID(id primitive.ObjectID) (Message, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	var message Message
	err := r.FindOne(ctx, bson.M{"_id": id}).Decode(&message)

	return message, err
}

func (r *MessagesRepo) InsertNewMessage(m Message) (Message, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	_, err := r.InsertOne(ctx, m)

	return m, err
}

// this function creates new ID and time and insert the document to database
func (r *MessagesRepo) InsertNewRawMessage(m Message) (Message, error) {
	m.ID = primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	m.CreatedAt = now
	m.UpdatedAt = now

	return r.InsertNewMessage(m)
}

func (r *MessagesRepo) GetMessagesOfConversation(
	conversationID primitive.ObjectID, limit int64,
) (*[]Message, error) {
	ctx, cal := context.WithTimeout(context.Background(), time.Second)
	defer cal()

	filter := bson.M{"conversationId": conversationID}
	messages := make([]Message, 0)
	cur, err := r.Find(ctx, filter,
		&options.FindOptions{Sort: bson.M{"createdAt": -1}, Limit: &limit})
	if err != nil {
		log.Println("can not get conversations:", err)
		return nil, err
	}
	err = cur.All(ctx, &messages)
	if err != nil {
		log.Println("can not parse conversations:", err)
		return nil, err
	}

	return &messages, nil
}
