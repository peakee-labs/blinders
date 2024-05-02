// Event represents an async event
package transport

import (
	"time"
)

type EventType string

/*
 * Transport interface of user/friends service
 */
const (
	AddFriend EventType = "ADD_FRIEND"
)

type Event struct {
	Type      EventType `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload"`
}

type AddFriendAction string

const (
	InitFriendRequest   AddFriendAction = "INIT_FRIEND_REQUEST"
	AcceptFriendRequest AddFriendAction = "ACCEPT_FRIEND_REQUEST"
	DenyFriendRequest   AddFriendAction = "DENY_FRIEND_REQUEST"
)

type AddFriendEvent struct {
	Event              `                json:",inline"`
	Action             AddFriendAction
	UserID             string `json:"userId"`
	AddFriendRequestID string `json:"addFriendRequestId"`
}

/*
 * Transport interface of collecting service
 */
const (
	AddTranslateLog EventType = "ADD_TRANSLATE_LOG"
	AddExplainLog   EventType = "ADD_EXPLAIN_LOG"
)
