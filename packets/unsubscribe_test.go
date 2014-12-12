package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUnsubscribeStruct(t *testing.T) {
	msg := NewMessage(UNSUBSCRIBE).(*UnsubscribeMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.UnsubscribeMessage", reflect.TypeOf(msg).String(), "Type should be UnsubscribeMessage")
		assert.Equal(t, 0, msg.TopicIdType, "Default TopicIdType should be 0")
		assert.Equal(t, 0, msg.TopicId, "Default TopicId should be 0")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")
		assert.Equal(t, []byte(nil), msg.TopicName, "Default TopicName should be blank")

		assert.Equal(t, UNSUBSCRIBE, msg.MessageType(), "MessageType() should return UNSUBSCRIBE")
	}
}
