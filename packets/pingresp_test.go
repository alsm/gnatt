package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPingrespStruct(t *testing.T) {
	msg := NewMessage(PINGRESP).(*PingrespMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.PingrespMessage", reflect.TypeOf(msg).String(), "Type should be PingrespMessage")
		assert.Equal(t, 2, msg.Length, "Default Length should be 2")

		assert.Equal(t, PINGRESP, msg.MessageType(), "MessageType() should return PINGRESP")
	}
}
