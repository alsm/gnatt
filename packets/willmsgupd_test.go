package packets

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillMsgUpdateStruct(t *testing.T) {
	msg := NewMessage(WILLMSGUPD).(*WillMsgUpdateMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.WillMsgUpdateMessage", reflect.TypeOf(msg).String(), "Type should be WillMsgUpdateMessage")
		assert.Equal(t, []byte(nil), msg.WillMsg, "Default WillMsg should be blank")

		assert.Equal(t, WILLMSGUPD, msg.MessageType(), "MessageType() should return WILLMSGUPD")
	}
}
