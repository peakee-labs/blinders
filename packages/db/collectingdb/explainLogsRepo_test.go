package collectingdb

import (
	"testing"
	"time"

	dbutils "blinders/packages/db/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	eclient, _     = dbutils.InitMongoClient("mongodb://localhost:27017")
	explainLogRepo = NewExplainLogsRepo(eclient.Database("blinders"))
)

func TestGetExplainLogWithSmallestGetCount(t *testing.T) {
	userID := primitive.NewObjectID()
	_, err := explainLogRepo.InsertRaw(&ExplainLog{
		UserID:   userID,
		GetCount: 1,
		Request:  ExplainRequest{Text: "hello1"},
	})
	assert.Nil(t, err)
	_, err = explainLogRepo.InsertRaw(&ExplainLog{
		UserID:   userID,
		GetCount: 0,
		Request:  ExplainRequest{Text: "hello2"},
	})
	assert.Nil(t, err)

	log, err := explainLogRepo.GetLogWithSmallestGetCountByUserID(userID)
	assert.Nil(t, err)
	assert.Equal(t, log.Request.Text, "hello2")
	assert.Equal(t, log.GetCount, 1)
}
func TestGetExplainLogWithPagination(t *testing.T) {
	userID := primitive.NewObjectID()
	logs := []*ExplainLog{
		{
			UserID:   userID,
			GetCount: 1,
			Request:  ExplainRequest{Text: "hello1"},
		},
		{
			UserID:   userID,
			GetCount: 0,
			Request:  ExplainRequest{Text: "hello2"},
		},
		{
			UserID:   userID,
			GetCount: 0,
			Request:  ExplainRequest{Text: "hello3"},
		},
	}

	insertedLogs := make([]*ExplainLog, len(logs))

	for idx, log := range logs {
		insertedLog, err := explainLogRepo.InsertRaw(log)
		assert.Nil(t, err)
		time.Sleep(10 * time.Millisecond)
		insertedLogs[idx] = insertedLog
	}

	// get 1 log
	opt := Pagination{
		Limit: 1,
		From:  time.Time{},
		To:    time.Now(),
	}

	logs, newPag, err := explainLogRepo.GetLogWithPagination(userID, &opt)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(logs))
	assert.Equal(t, insertedLogs[0].Request, logs[0].Request)
	assert.Equal(t, insertedLogs[0].Response, logs[0].Response)
	assert.Equal(t, insertedLogs[0].CreatedAt.Time(), newPag.From)
	assert.Equal(t, insertedLogs[0].CreatedAt.Time(), newPag.To)

	// get all log
	opt = Pagination{
		Limit: len(insertedLogs) + 1,
		From:  time.Time{},
		To:    time.Now(),
	}

	logs, newPag, err = explainLogRepo.GetLogWithPagination(userID, &opt)
	assert.Nil(t, err)
	assert.Equal(t, len(insertedLogs), len(logs))

	for idx, log := range logs {
		assert.Equal(t, insertedLogs[idx].Request, log.Request)
		assert.Equal(t, insertedLogs[idx].Response, log.Response)
	}
	assert.Equal(t, insertedLogs[0].CreatedAt.Time(), newPag.From)
	assert.Equal(t, insertedLogs[len(insertedLogs)-1].CreatedAt.Time(), newPag.To)

	// get logs from the second log
	fromSecond := insertedLogs[0].CreatedAt.Time()
	opt = Pagination{
		Limit: len(insertedLogs) + 1,
		From:  fromSecond,
		To:    time.Now(),
	}
	logs, newPag, err = explainLogRepo.GetLogWithPagination(userID, &opt)
	assert.Nil(t, err)
	assert.Equal(t, len(insertedLogs)-1, len(logs))

	for idx, log := range logs {
		assert.Equal(t, insertedLogs[idx+1].Request, log.Request)
		assert.Equal(t, insertedLogs[idx+1].Response, log.Response)
	}
	assert.Equal(t, insertedLogs[1].CreatedAt.Time(), newPag.From)
	assert.Equal(t, insertedLogs[len(insertedLogs)-1].CreatedAt.Time(), newPag.To)

	// get with empty pagination
	logs, newPag, err = explainLogRepo.GetLogWithPagination(userID, nil)
	assert.Nil(t, err)
	assert.Equal(t, len(insertedLogs), len(logs))

	for idx, log := range logs {
		assert.Equal(t, insertedLogs[idx].Request, log.Request)
		assert.Equal(t, insertedLogs[idx].Response, log.Response)
	}
	assert.Equal(t, insertedLogs[0].CreatedAt.Time(), newPag.From)
	assert.Equal(t, insertedLogs[len(insertedLogs)-1].CreatedAt.Time(), newPag.To)
}
