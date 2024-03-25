package wschat

import "blinders/packages/db/models"

type ChatEventType string

const (
	UserPing                  ChatEventType = "USER:PING"
	UserSendMessage           ChatEventType = "USER:SEND_MESSAGE"
	UserUpdateMessageStatus   ChatEventType = "USER:UPDATE_MESSAGE_STATUS"
	ServerSendMessage         ChatEventType = "SERVER:SEND_MESSAGE"
	ServerAckSendMessage      ChatEventType = "SERVER:ACK_SEND_MESSAGE"
	ServerUpdateMessageStatus ChatEventType = "SERVER:UPDATE_MESSAGE_STATUS"
)

type ChatEvent struct {
	Type ChatEventType `json:"type"`
}

type AckError struct {
	Error string `json:"error"`
}

type UserSendMessagePayload struct {
	ChatEvent      `json:",inline"`
	Content        string `json:"content"`
	ConversationID string `json:"conversationId"`
	ReplyTo        string `json:"replyTo"`
	ResolveID      string `json:"resolveId"` // it helps client side resolve the message
}

type ServerAckSendMessagePayload struct {
	ChatEvent `json:",inline"`
	ResolveID string         `json:"resolveId"` // send ack response to sender
	Message   models.Message `json:"message,omitempty"`
	Error     AckError       `json:"error,omitempty"`
}

type ServerSendMessagePayload struct {
	ChatEvent `json:",inline"`
	Message   models.Message `json:"message"`
}

type MessageStatus string

type UserUpdateMessageStatusPayload struct {
	ChatEvent      `json:",inline"`
	Content        string `json:"content"`
	ConversationID string `json:"conversationId"`
	MessageID      string `json:"messageId"`
	Status         string
}
