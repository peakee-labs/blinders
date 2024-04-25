package collecting_test

import (
	"blinders/packages/collecting"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testURL = "mongodb://localhost:27017"
	testDB  = "blinders"
)

func Test_TranslateLog(t *testing.T) {
	db := GetDatabase(t)

	eventCollector := collecting.NewEventCollector(db)
	defer eventCollector.Col.Drop(context.Background())

	userOID := primitive.NewObjectID()

	getLog, err := eventCollector.GetTranslateLogByUserID(userOID)
	assert.Nil(t, err)
	assert.Empty(t, getLog)

	translateLog := &collecting.TranslateEvent{
		UserID: userOID,
		Request: collecting.TranslateRequest{
			Text: "text",
		},
		Response: collecting.TranslateResponse{
			Translate: "translate",
		},
	}
	addedLog, err := eventCollector.AddRawTranslateLog(&collecting.TranslateEventLog{
		TranslateEvent: *translateLog,
	})
	assert.NoError(t, err)
	assert.NotNil(t, addedLog)

	getLog, err = eventCollector.GetTranslateLogByUserID(userOID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getLog)

	assert.Equal(t, 1, len(getLog))
	assert.Equal(t, *translateLog, getLog[0].TranslateEvent)

	anotherTranslateLog := &collecting.TranslateEvent{
		UserID: userOID,
		Request: collecting.TranslateRequest{
			Text: "anotherText",
		},
		Response: collecting.TranslateResponse{
			Translate: "anotherTranslate",
		},
	}
	anotherAddedLog, err := eventCollector.AddRawTranslateLog(&collecting.TranslateEventLog{
		TranslateEvent: *anotherTranslateLog,
	})
	assert.NoError(t, err)
	assert.NotNil(t, anotherAddedLog)

	getLog, err = eventCollector.GetTranslateLogByUserID(userOID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getLog)

	assert.Equal(t, 2, len(getLog))

	l, err := eventCollector.GetTranslateLogByID(anotherAddedLog.ID)
	assert.NoError(t, err)
	assert.Equal(t, *anotherAddedLog, *l)
}

func Test_ExplainLog(t *testing.T) {
	db := GetDatabase(t)

	eventCollector := collecting.NewEventCollector(db)
	defer eventCollector.Col.Drop(context.Background())

	userOID := primitive.NewObjectID()

	getLog, err := eventCollector.GetExplainLogByUserID(userOID)
	assert.Nil(t, err)
	assert.Empty(t, getLog)

	explainLog := &collecting.ExplainEvent{
		UserID: userOID,
		Request: collecting.ExplainRequest{
			Text:     "text",
			Sentence: "sentence",
		},
		Response: collecting.ExplainResponse{
			Translate:       "translate",
			ExpandWords:     make([]string, 0),
			KeyWords:        make([]string, 0),
			GrammarAnalysis: make(map[string]any, 0),
		},
	}
	addedLog, err := eventCollector.AddRawExplainLog(&collecting.ExplainEventLog{
		ExplainEvent: *explainLog,
	})
	assert.NoError(t, err)
	assert.NotNil(t, addedLog)

	getLog, err = eventCollector.GetExplainLogByUserID(userOID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getLog)

	assert.Equal(t, 1, len(getLog))
	assert.Equal(t, *explainLog, getLog[0].ExplainEvent)

	l, err := eventCollector.GetExplainLogByID(addedLog.ID)
	assert.NoError(t, err)
	assert.Equal(t, *addedLog, *l)

	anotherExplainLog := &collecting.ExplainEvent{
		UserID: userOID,
		Request: collecting.ExplainRequest{
			Text:     "anothertext",
			Sentence: "anothersentence",
		},
		Response: collecting.ExplainResponse{
			Translate:       "anothertranslate",
			ExpandWords:     make([]string, 0),
			KeyWords:        make([]string, 0),
			GrammarAnalysis: make(map[string]any, 0),
		},
	}
	anotherAddedLog, err := eventCollector.AddRawExplainLog(&collecting.ExplainEventLog{
		ExplainEvent: *anotherExplainLog,
	})
	assert.NoError(t, err)
	assert.NotNil(t, anotherAddedLog)

	getLog, err = eventCollector.GetExplainLogByUserID(userOID)
	assert.NoError(t, err)
	assert.NotEmpty(t, getLog)

	assert.Equal(t, 2, len(getLog))
}

func GetDatabase(t *testing.T) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(testURL))
	if err != nil {
		t.Fatalf("Error connecting to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		t.Fatalf("Error ping to MongoDB: %v", err)
	}

	// Create database and collection
	db := client.Database(testDB)
	return db
}
