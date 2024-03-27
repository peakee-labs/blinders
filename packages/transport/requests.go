package transport

import "blinders/packages/db/models"

type RequestType string

const (
	Embedding        RequestType = "EMBEDDING"
	AddUserMatchInfo RequestType = "ADD_USER_MATCH_INFO"
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
