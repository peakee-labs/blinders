package wschat

import (
	"testing"

	dbutils "blinders/packages/dbutils"
	"blinders/packages/session"
	chatrepo "blinders/services/chat/repo"
	usersrepo "blinders/services/users/repo"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userRepo *usersrepo.UsersRepo

func init() {
	client, _ := dbutils.InitMongoClient("mongodb://localhost:27017")
	InitChatApp(
		session.NewManager(redis.NewClient(&redis.Options{Addr: "localhost:6379"})),
		client.Database("blinders"),
	)
	userRepo = usersrepo.NewUsersRepo(client.Database("blinders"))
}

func TestSendMessageFailedWithWrongPayload(t *testing.T) {
	_, err := HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: "wrongID",
			ResolveID:      "resolveID",
		})

	assert.NotNil(t, err)

	_, err = HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: primitive.NewObjectID().Hex(),
			ResolveID:      "resolveID",
		})

	assert.NotNil(t, err)
}

func TestSendMessageFailedWithConversationNotFound(t *testing.T) {
	_, err := HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: primitive.NewObjectID().Hex(),
		})

	assert.NotNil(t, err)
}

func TestSendMessageFailedWithUserIsNotMember(t *testing.T) {
	conversation, _ := app.ConvsRepo.InsertNewConversation(chatrepo.Conversation{})
	_, err := HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: conversation.ID.Hex(),
		})

	assert.NotNil(t, err)
}

func TestSendMessageWithNoError(t *testing.T) {
	user, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	conversation, _ := app.ConvsRepo.InsertNewRawConversation(chatrepo.Conversation{
		Members: []chatrepo.Member{{UserID: user.ID}},
	})
	_, err := HandleSendMessage(
		user.ID.Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: conversation.ID.Hex(),
		})

	assert.Nil(t, err)
}

func TestSendMessageFailedWithInvalidMessageToReply(t *testing.T) {
	_, err := HandleSendMessage(
		primitive.NewObjectID().Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: primitive.NewObjectID().Hex(),
			ReplyTo:        primitive.NewObjectID().Hex(),
		})

	assert.NotNil(t, err)
}

func TestSendMessageWithValidMessageToReply(t *testing.T) {
	user, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	conversation, _ := app.ConvsRepo.InsertNewRawConversation(chatrepo.Conversation{
		Members: []chatrepo.Member{{UserID: user.ID}},
	})
	message, _ := app.MessagesRepo.InsertNewRawMessage(chatrepo.Message{
		ConversationID: conversation.ID,
	})

	_, err := HandleSendMessage(
		user.ID.Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: conversation.ID.Hex(),
			ReplyTo:        message.ID.Hex(),
		})

	assert.Nil(t, err)
}

func TestSendMessageSuccess(t *testing.T) {
	user, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	conversation, _ := app.ConvsRepo.InsertNewRawConversation(chatrepo.Conversation{
		Members: []chatrepo.Member{{UserID: user.ID}},
	})
	_, err := HandleSendMessage(
		user.ID.Hex(),
		primitive.NewObjectID().Hex(),
		UserSendMessagePayload{
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        "hello world",
			ConversationID: conversation.ID.Hex(),
		})

	assert.Nil(t, err)
}

func TestSendMessageWithDistribution(t *testing.T) {
	sender, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	recipient1, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	recipient2, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	conversation, _ := app.ConvsRepo.InsertNewRawConversation(
		chatrepo.Conversation{
			Members: []chatrepo.Member{
				{UserID: sender.ID},
				{UserID: recipient1.ID},
				{UserID: recipient2.ID},
			},
		})

	sConnID := primitive.NewObjectID().Hex()
	r1connID := primitive.NewObjectID().Hex()
	r2connID := primitive.NewObjectID().Hex()
	_ = app.Session.AddSession(recipient1.ID.Hex(), r1connID)
	_ = app.Session.AddSession(recipient2.ID.Hex(), r2connID)

	resolveID := primitive.NewObjectID().Hex()
	content := "hello world"
	dCh, err := HandleSendMessage(
		sender.ID.Hex(),
		sConnID,
		UserSendMessagePayload{
			ResolveID:      resolveID,
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        content,
			ConversationID: conversation.ID.Hex(),
		})

	assert.Nil(t, err)

	expectedMap := map[string]bool{}
	for {
		de := <-dCh
		if de == nil {
			break
		}
		expectedMap[de.ConnectionID] = true
		switch de.ConnectionID {
		case sConnID:
			payload := de.Payload.(ServerAckSendMessagePayload)
			assert.Equal(t, ServerAckSendMessage, payload.Type)
			assert.Equal(t, conversation.ID, payload.Message.ConversationID)
			assert.Equal(t, "", payload.Error.Error)
			assert.Equal(t, content, payload.Message.Content)
			assert.Equal(t, resolveID, payload.ResolveID)
		case r1connID:
			payload := de.Payload.(ServerSendMessagePayload)
			assert.Equal(t, ServerSendMessage, payload.Type)
			assert.Equal(t, conversation.ID, payload.Message.ConversationID)
			assert.Equal(t, content, payload.Message.Content)
		case r2connID:
			payload := de.Payload.(ServerSendMessagePayload)
			assert.Equal(t, ServerSendMessage, payload.Type)
			assert.Equal(t, conversation.ID, payload.Message.ConversationID)
			assert.Equal(t, content, payload.Message.Content)
		}

	}

	assert.True(t, expectedMap[r1connID])
	assert.True(t, expectedMap[r2connID])
	assert.True(t, expectedMap[sConnID])
	assert.Equal(t, 3, len(expectedMap))
}

func TestSendMessageWithDistributionWithOfflineRecipient(t *testing.T) {
	sender, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	recipient1, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	recipient2, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	conversation, _ := app.ConvsRepo.InsertNewRawConversation(
		chatrepo.Conversation{
			Members: []chatrepo.Member{
				{UserID: sender.ID},
				{UserID: recipient1.ID},
				{UserID: recipient2.ID},
			},
		})

	sConnID := primitive.NewObjectID().Hex()
	r1connID := primitive.NewObjectID().Hex()
	_ = app.Session.AddSession(recipient1.ID.Hex(), r1connID)

	resolveID := primitive.NewObjectID().Hex()
	content := "hello world"
	dCh, err := HandleSendMessage(
		sender.ID.Hex(),
		sConnID,
		UserSendMessagePayload{
			ResolveID:      resolveID,
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        content,
			ConversationID: conversation.ID.Hex(),
		})

	assert.Nil(t, err)

	expectedMap := map[string]bool{}
	for {
		de := <-dCh
		if de == nil {
			break
		}
		expectedMap[de.ConnectionID] = true

	}

	assert.True(t, expectedMap[r1connID])
	assert.True(t, expectedMap[sConnID])
	assert.Equal(t, 2, len(expectedMap))
}

func TestSendMessageWithDistributionWithMultipleSessionsPerUser(t *testing.T) {
	sender, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	recipient1, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	recipient2, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	conversation, _ := app.ConvsRepo.InsertNewRawConversation(
		chatrepo.Conversation{
			Members: []chatrepo.Member{
				{UserID: sender.ID},
				{UserID: recipient1.ID},
				{UserID: recipient2.ID},
			},
		})

	sConnID := primitive.NewObjectID().Hex()
	r1connID := primitive.NewObjectID().Hex()
	r1connID2 := primitive.NewObjectID().Hex()
	_ = app.Session.AddSession(recipient1.ID.Hex(), r1connID)
	_ = app.Session.AddSession(recipient1.ID.Hex(), r1connID2)

	resolveID := primitive.NewObjectID().Hex()
	content := "hello world"
	dCh, err := HandleSendMessage(
		sender.ID.Hex(),
		sConnID,
		UserSendMessagePayload{
			ResolveID:      resolveID,
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        content,
			ConversationID: conversation.ID.Hex(),
		})

	assert.Nil(t, err)

	expectedMap := map[string]bool{}
	for {
		de := <-dCh
		if de == nil {
			break
		}
		expectedMap[de.ConnectionID] = true

	}

	assert.True(t, expectedMap[r1connID])
	assert.True(t, expectedMap[r1connID2])
	assert.True(t, expectedMap[sConnID])
	assert.Equal(t, 3, len(expectedMap))
}

func TestSendMessageWithDistributionWithStoredMessage(t *testing.T) {
	sender, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	recipient1, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	recipient2, _ := userRepo.InsertNewRawUser(usersrepo.User{})
	conversation, _ := app.ConvsRepo.InsertNewRawConversation(
		chatrepo.Conversation{
			Members: []chatrepo.Member{
				{UserID: sender.ID},
				{UserID: recipient1.ID},
				{UserID: recipient2.ID},
			},
		})

	sConnID := primitive.NewObjectID().Hex()
	r1connID := primitive.NewObjectID().Hex()
	_ = app.Session.AddSession(recipient1.ID.Hex(), r1connID)

	resolveID := primitive.NewObjectID().Hex()
	content := "hello world"
	dCh, err := HandleSendMessage(
		sender.ID.Hex(),
		sConnID,
		UserSendMessagePayload{
			ResolveID:      resolveID,
			ChatEvent:      ChatEvent{Type: UserSendMessage},
			Content:        content,
			ConversationID: conversation.ID.Hex(),
		})

	assert.Nil(t, err)

	var message1 chatrepo.Message
	var message2 chatrepo.Message
	for {
		de := <-dCh
		if de == nil {
			break
		}
		if de.ConnectionID == sConnID {
			message1 = de.Payload.(ServerAckSendMessagePayload).Message
		} else if de.ConnectionID == r1connID {
			message2 = de.Payload.(ServerSendMessagePayload).Message
		}

	}
	storedMessage, err := app.MessagesRepo.GetMessageByID(message1.ID)
	assert.Nil(t, err)
	assert.Equal(t, storedMessage, message1)
	assert.Equal(t, storedMessage, message2)
}
