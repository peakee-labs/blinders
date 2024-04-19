package transport

import "time"

type EventType string

const (
	AddFriend EventType = "ADD_FRIEND"
)

type Event struct {
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

type AddFriendAction string

const (
	InitFriendRequest   AddFriendAction = "INIT"
	AcceptFriendRequest AddFriendAction = "ACCEPT"
	DenyFriendRequest   AddFriendAction = "DENY"
)

type AddFriendEvent struct {
	Event              `json:",inline"`
	UserID             string `json:"userId"`
	AddFriendRequestID string `json:"addFriendRequestId"`
	Action             AddFriendAction
}

// GenericEvent could be used to embed specific event as payload.
//
// The event producer and consumer could identify payload type by checking
// Event.Type value of GenericEvent.Event field
type GenericEvent struct {
	Event   `json:",inline"`
	Payload any `json:",inline"`
}
