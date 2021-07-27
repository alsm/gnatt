package packets

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingreqStruct(t *testing.T) {
	msg := NewMessage(PINGREQ).(*PingreqMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.PingreqMessage", reflect.TypeOf(msg).String(), "Type should be PingreqMessage")
		assert.Equal(t, []byte(nil), msg.ClientId, "Default ClientId should be blank")

		assert.Equal(t, PINGREQ, msg.MessageType(), "MessageType() should return PINGREQ")
	}
}
