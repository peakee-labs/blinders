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
