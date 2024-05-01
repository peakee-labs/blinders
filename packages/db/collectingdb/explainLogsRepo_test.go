package collectingdb

import (
	"testing"

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
