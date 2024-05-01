// Request represents a sync request
package transport

import (
	"blinders/packages/db/collectingdb"
	"blinders/packages/db/matchingdb"
)

type RequestType string

const (
	Embedding        RequestType = "EMBEDDING"
	AddUserMatchInfo RequestType = "ADD_USER_MATCH_INFO"
	GetTranslateLog  RequestType = "GET_TRANSLATE_LOG"
	GetExplainLog    RequestType = "GET_EXPLAIN_LOG"
)

type Request struct {
	Type RequestType `json:"type"`
}

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
	Request `       json:",inline"`
	UserID  string `json:"userId"`
}

type GetTranslateLogResponse struct {
	collectingdb.TranslateLog
}

type GetExplainLogResponse struct {
	collectingdb.ExplainLog
}
