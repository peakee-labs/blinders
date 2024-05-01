package collectingdb

import "go.mongodb.org/mongo-driver/mongo"

var (
	ExplainLogsCollection   = "explain-logs"
	TranslateLogsCollection = "translate-logs"
)

type CollectingDB struct {
	ExplainLogsRepo   *ExplainLogsRepo
	TranslateLogsRepo *TranslateLogsRepo
}

func NewCollectingDB(db *mongo.Database) *CollectingDB {
	return &CollectingDB{
		ExplainLogsRepo:   NewExplainLogsRepo(db),
		TranslateLogsRepo: NewTranslateLogsRepo(db),
	}
}
