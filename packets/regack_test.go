package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRegackStruct(t *testing.T) {
	msg := NewMessage(REGACK).(*RegackMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.RegackMessage", reflect.TypeOf(msg).String(), "Type should be RegackMessage")
		assert.Equal(t, 7, msg.Length, "Default Length should be 7")
		assert.Equal(t, 0, msg.TopicId, "Default TopicId should be 0")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")
		assert.Equal(t, 0, msg.ReturnCode, "Default ReturnCode should be 0")

		assert.Equal(t, REGACK, msg.MessageType(), "MessageType() should return REGACK")
	}
}
