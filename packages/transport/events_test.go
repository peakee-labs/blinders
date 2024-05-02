package transport

import (
	"testing"

	"blinders/packages/db/collectingdb"
	"blinders/packages/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestParseGenericEventWithInlineAnyViaAny(t *testing.T) {
	originalEvent := AddTranslateLogEvent{
		Event{Type: AddTranslateLog},
		collectingdb.TranslateLog{
			UserID: primitive.NewObjectID(),
		},
	}

	anyEvent, _ := utils.JSONConvert[any](originalEvent)
	genericEvent, _ := utils.JSONConvert[Event](anyEvent)
	targetEvent, _ := utils.JSONConvert[AddTranslateLogEvent](genericEvent)
	assert.Equal(t, targetEvent.Payload.UserID, originalEvent.Payload.UserID)
}

func TestParseGenericEventFailedWithInlineAnyViaEvent(t *testing.T) {
	originalEvent := AddTranslateLogEvent{
		Event{Type: AddTranslateLog},
		collectingdb.TranslateLog{
			UserID: primitive.NewObjectID(),
		},
	}

	genericEvent, _ := utils.JSONConvert[Event](originalEvent)
	targetEvent, _ := utils.JSONConvert[AddTranslateLogEvent](genericEvent)
	assert.Equal(t, targetEvent.Payload.UserID, originalEvent.Payload.UserID)
}
