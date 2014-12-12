package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGwInfoStruct(t *testing.T) {
	msg := NewMessage(GWINFO).(*GwInfoMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.GwInfoMessage", reflect.TypeOf(msg).String(), "Type should be GwInfoMessage")
		assert.Equal(t, 0, msg.GatewayId, "Default GatewayId should be 0")
		assert.Equal(t, []byte(nil), msg.GatewayAddress, "Default GatewayAddress should be blank")

		assert.Equal(t, GWINFO, msg.MessageType(), "MessageType() should return GWINFO")
	}
}
