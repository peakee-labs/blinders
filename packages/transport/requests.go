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
		// this inline any keeps all the fields of the original request,
		// then converting to actual request by type
		Any any `json:",inline"`
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
	Request `       json:",inline"`
	Data    string
}
type EmbeddingResponse struct {
	Embedded []float32
}

/*
 * Transport interface of explore service
 */
type AddUserMatchInfoRequest struct {
	Request `json:",inline"`
	Data    matchingdb.MatchInfo `json:"data"`
}
type AddUserMatchInfoResponse struct {
	Error *string `json:"error,omitempty"`
}

/*
 * Transport interface of collecting service
 */
type GetCollectingLogRequest struct {
	Request `json:",inline"`
	UserID  string `json:"userId"`
}

type GetTranslateLogResponse struct {
	collectingdb.TranslateLog
}

type GetExplainLogResponse struct {
	collectingdb.ExplainLog
}
