package suggestapi

import (
	"fmt"
	"reflect"
	"time"

	"blinders/packages/suggest"
)

func (m ClientMessage) ToCommonMessage() suggest.Message {
	var Timestamp int64
	switch timestamp := m.Timestamp.(type) {
	case int:
		Timestamp = int64(timestamp)
	case string:
		// expect date time as string type, "Tue Dec 05 2023 12:35:04 GMT+0700"
		layout := "Mon Jan 02 2006 15:04:05 GMT-0700"
		t, err := time.Parse(layout, timestamp)
		if err != nil {
			panic(fmt.Sprintf("clientMessage: given time (%s) cannot parse with layout (%s)", timestamp, layout))
		}
		Timestamp = t.Unix()
	default:
		panic(fmt.Sprintf("clientMessage: unknown timestamp type (%s)", reflect.TypeOf(m.Timestamp).String()))
	}

	return suggest.Message{
		Sender:    m.Sender,
		Receiver:  m.Receiver,
		Content:   m.Content,
		Timestamp: Timestamp,
	}
}
