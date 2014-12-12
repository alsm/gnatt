package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSubscribeStruct(t *testing.T) {
	msg := NewMessage(SUBSCRIBE).(*SubscribeMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.SubscribeMessage", reflect.TypeOf(msg).String(), "Type should be SubscribeMessage")
		assert.Equal(t, false, msg.Dup, "Default Dup flag should be false")
		assert.Equal(t, 0, msg.Qos, "Default Qos should be 0")
		assert.Equal(t, 0, msg.TopicIdType, "Default TopicIdType should be 0")
		assert.Equal(t, 0, msg.TopicId, "Default TopicId should be 0")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")
		assert.Equal(t, []byte(nil), msg.TopicName, "Default Topicname should be blank")

		assert.Equal(t, SUBSCRIBE, msg.MessageType(), "MessageType() should return SUBSCRIBE")
	}
}
