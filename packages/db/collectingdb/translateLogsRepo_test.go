package collectingdb

import (
	"testing"

	dbutils "blinders/packages/db/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	client, _        = dbutils.InitMongoClient("mongodb://localhost:27017")
	translateLogRepo = NewTranslateLogsRepo(client.Database("blinders"))
)

func TestGetTranslateLogWithSmallestGetCount(t *testing.T) {
	userID := primitive.NewObjectID()
	_, err := translateLogRepo.InsertRaw(&TranslateLog{
		UserID:   userID,
		GetCount: 1,
		Request:  TranslateRequest{Text: "hello1"},
	})
	assert.Nil(t, err)
	_, err = translateLogRepo.InsertRaw(&TranslateLog{
		UserID:   userID,
		GetCount: 0,
		Request:  TranslateRequest{Text: "hello2"},
	})
	assert.Nil(t, err)

	log, err := translateLogRepo.GetLogWithSmallestGetCountByUserID(userID)
	assert.Nil(t, err)
	assert.Equal(t, log.Request.Text, "hello2")
	assert.Equal(t, log.GetCount, 1)
}
