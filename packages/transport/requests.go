// Request represents a sync request
package transport

import (
	"blinders/packages/db/collectingdb"
	"blinders/packages/db/matchingdb"
)

type (
	RequestType string

	// this Request struct is the base struct for all requests
	// do not use it directly
	Request struct {
		Type RequestType `json:"type"`
		// This payload field keeps the original payload,
		// all the event must put their payload in this field to prevent missing fields
		Payload any `json:"payload"`
	}
)

const (
	Embedding        RequestType = "EMBEDDING"
	AddUserMatchInfo RequestType = "ADD_USER_MATCH_INFO"
	GetTranslateLog  RequestType = "GET_TRANSLATE_LOG"
	GetExplainLog    RequestType = "GET_EXPLAIN_LOG"
)

/*
 * For vector embedding
 */
type EmbeddingRequest struct {
	Request `json:",inline"`
	// This payload field keeps the original payload,
	// all the event must put their payload in this field to prevent missing fields
	Payload string `json:"payload"`
}
type EmbeddingResponse struct {
	Embedded []float32
}

/*
 * Transport interface of explore service
 */
type AddUserMatchInfoRequest struct {
	Request `json:",inline"`
	Payload matchingdb.MatchInfo `json:"payload"`
}
type AddUserMatchInfoResponse struct {
	Error *string `json:"error,omitempty"`
}

/*
 * Transport interface of collecting service
 */
type GetCollectingLogRequest struct {
	Request `json:",inline"`
	Payload GetCollectingLogPayload `json:"payload"`
}

type GetCollectingLogPayload struct {
	UserID string `json:"userId"`
}

type GetTranslateLogResponse struct {
	collectingdb.TranslateLog
}

type GetExplainLogResponse struct {
	collectingdb.ExplainLog
}
