package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPubrelStruct(t *testing.T) {
	msg := NewMessage(PUBREL).(*PubrelMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.PubrelMessage", reflect.TypeOf(msg).String(), "Type should be PubrelMessage")
		assert.Equal(t, 4, msg.Length, "Default Length should be 4")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")

		assert.Equal(t, PUBREL, msg.MessageType(), "MessageType() should return PUBREL")
	}
}
