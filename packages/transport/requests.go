package transport

import (
	"blinders/packages/collecting"
	"blinders/packages/db/models"
)

type RequestType string

const (
	Embedding        RequestType = "EMBEDDING"
	AddUserMatchInfo RequestType = "ADD_USER_MATCH_INFO"
	CollectEvent     RequestType = "COLLECT_EVENT"
	GetEvent         RequestType = "GET_EVENT"
)

type Request struct {
	Type RequestType `json:"type"`
}

type EmbeddingRequest struct {
	Request `json:",inline"`
	Data    string
}

type EmbeddingResponse struct {
	Embedded []float32
}

type AddUserMatchInfoRequest struct {
	Request `json:",inline"`
	Data    models.MatchInfo `json:"data"`
}

type AddUserMatchInfoResponse struct {
	Error *string `json:"error,omitempty"`
}

type CollectEventRequest struct {
	Request `json:",inline"`
	Data    collecting.GenericEvent `json:"data"`
}

type GetEventRequest struct {
	Request   `json:",inline"`
	UserID    string               `json:"userId"`
	Type      collecting.EventType `json:"eventType"`
	NumReturn int                  `json:"numReturn"`
}

type GetEventResponse struct {
	Data []collecting.GenericEvent `json:"data"`
}
