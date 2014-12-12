package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPubrecStruct(t *testing.T) {
	msg := NewMessage(PUBREC).(*PubrecMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.PubrecMessage", reflect.TypeOf(msg).String(), "Type should be PubrecMessage")
		assert.Equal(t, 4, msg.Length, "Default Length should be 4")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")

		assert.Equal(t, PUBREC, msg.MessageType(), "MessageType() should return PUBREC")
	}
}
