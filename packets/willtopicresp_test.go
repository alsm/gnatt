package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestWillTopicRespStruct(t *testing.T) {
	msg := NewMessage(WILLTOPICRESP).(*WillTopicRespMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.WillTopicRespMessage", reflect.TypeOf(msg).String(), "Type should be WillTopicRespMessage")
		assert.Equal(t, 0, msg.ReturnCode, "Default ReturnCode should be 0")

		assert.Equal(t, WILLTOPICRESP, msg.MessageType(), "MessageType() should return WILLTOPICRESP")
	}
}
