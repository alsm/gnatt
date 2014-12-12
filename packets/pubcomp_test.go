package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPubcompStruct(t *testing.T) {
	msg := NewMessage(PUBCOMP).(*PubcompMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.PubcompMessage", reflect.TypeOf(msg).String(), "Type should be PubcompMessage")
		assert.Equal(t, 4, msg.Length, "Default Length should be 4")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")

		assert.Equal(t, PUBCOMP, msg.MessageType(), "MessageType() should return PUBCOMP")
	}
}
