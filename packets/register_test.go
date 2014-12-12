package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRegisterStruct(t *testing.T) {
	msg := NewMessage(REGISTER).(*RegisterMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.RegisterMessage", reflect.TypeOf(msg).String(), "Type should be RegisterMessage")
		assert.Equal(t, 0, msg.TopicId, "Default TopicId should be 0")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")
		assert.Equal(t, []byte(nil), msg.TopicName, "Default TopicName should be blank")

		assert.Equal(t, REGISTER, msg.MessageType(), "MessageType() should return REGISTER")
	}
}
