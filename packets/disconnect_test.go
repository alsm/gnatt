package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDisconnectStruct(t *testing.T) {
	msg := NewMessage(DISCONNECT).(*DisconnectMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.DisconnectMessage", reflect.TypeOf(msg).String(), "Type should be DisconnectMessage")
		assert.Equal(t, 0, msg.Duration, "Default Duration should be 0")

		assert.Equal(t, DISCONNECT, msg.MessageType(), "MessageType() should return DISCONNECT")
	}
}
