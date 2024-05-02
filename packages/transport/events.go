// Event represents an async event
package transport

import (
	"time"

	"blinders/packages/db/collectingdb"
)

type (
	EventType string

	// this Event struct is the base struct for all events
	// do not use it directly
	Event struct {
		Type      EventType `json:"type"`
		Timestamp time.Time `json:"timestamp"`
		// This payload field keeps the original payload,
		// all the event must put their payload in this field to prevent missing fields
		Payload any `json:"payload"`
	}
)

/*
 * Transport interface of user/friends service
 */
const (
	AddFriend EventType = "ADD_FRIEND"
)

type AddFriendAction string

const (
	InitFriendRequest   AddFriendAction = "INIT_FRIEND_REQUEST"
	AcceptFriendRequest AddFriendAction = "ACCEPT_FRIEND_REQUEST"
	DenyFriendRequest   AddFriendAction = "DENY_FRIEND_REQUEST"
)

type AddFriendEvent struct {
	Event   `json:",inline"`
	Payload AddFriendPayload `json:"payload"`
}

type AddFriendPayload struct {
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

type AddTranslateLogEvent struct {
	Event   `json:",inline"`
	Payload collectingdb.TranslateLog `json:"payload"`
}

type AddExplainLogEvent struct {
	Event   `json:",inline"`
	Payload collectingdb.ExplainLog `json:"payload"`
}
