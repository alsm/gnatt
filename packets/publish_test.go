package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPublishStruct(t *testing.T) {
	msg := NewMessage(PUBLISH).(*PublishMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.PublishMessage", reflect.TypeOf(msg).String(), "Type should be PublishMessage")
		assert.Equal(t, false, msg.Dup, "Default Dup flag should be false")
		assert.Equal(t, false, msg.Retain, "Default Retain flag should be false")
		assert.Equal(t, 0, msg.Qos, "Default Qos should be 0")
		assert.Equal(t, 0, msg.TopicIdType, "Default TopicIdType should be 0")
		assert.Equal(t, 0, msg.TopicId, "Default TopicId should be 0")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")
		assert.Equal(t, []byte(nil), msg.Data, "Default Data should be blank")

		assert.Equal(t, PUBLISH, msg.MessageType(), "MessageType() should return PUBLISH")
	}
}
