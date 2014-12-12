package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestWillMsgStruct(t *testing.T) {
	msg := NewMessage(WILLMSG).(*WillMsgMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.WillMsgMessage", reflect.TypeOf(msg).String(), "Type should be WillMsgMessage")
		assert.Equal(t, []byte(nil), msg.WillMsg, "Default WillMsg should be blank")

		assert.Equal(t, WILLMSG, msg.MessageType(), "MessageType() should return WILLMSG")
	}
}
