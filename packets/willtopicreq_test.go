package packets

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillTopicReqStruct(t *testing.T) {
	msg := NewMessage(WILLTOPICREQ).(*WillTopicReqMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.WillTopicReqMessage", reflect.TypeOf(msg).String(), "Type should be WillTopicReqMessage")

		assert.Equal(t, WILLTOPICREQ, msg.MessageType(), "MessageType() should return WILLTOPICREQ")
	}
}
