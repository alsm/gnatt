package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestConnectMessage(t *testing.T) {
	msg := NewMessage(CONNECT).(*ConnectMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.ConnectMessage", reflect.TypeOf(msg).String(), "Type should be ConnectMessage")
		assert.Equal(t, false, msg.Will, "Default Will should be false")
		assert.Equal(t, false, msg.CleanSession, "Default CleanSession should be false")
		assert.Equal(t, 0x01, msg.ProtocolId, "Default ProtocolId should be 1")
		assert.Equal(t, 0, msg.Duration, "Default Duration should be 0")
		assert.Equal(t, []byte(nil), msg.ClientId, "Default ClientId should be blank")

		assert.Equal(t, CONNECT, msg.MessageType(), "MessageType() should reurn CONNECT")
	}
}
