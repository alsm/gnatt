package packets

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnackStruct(t *testing.T) {
	msg := NewMessage(CONNACK).(*ConnackMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.ConnackMessage", reflect.TypeOf(msg).String(), "Type should be ConnackMessage")
		assert.Equal(t, 0, msg.ReturnCode, "Default ReturnCode should be 0")
		assert.Equal(t, 3, msg.Length, "Length should be 3")

		assert.Equal(t, CONNACK, msg.MessageType(), "MessageType() should return CONNACK")
	}
}
