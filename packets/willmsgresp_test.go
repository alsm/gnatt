package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestWillMsgRespStruct(t *testing.T) {
	msg := NewMessage(WILLMSGRESP).(*WillMsgRespMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.WillMsgRespMessage", reflect.TypeOf(msg).String(), "Type should be WillMsgRespMessage")
		assert.Equal(t, 0, msg.ReturnCode, "Default ReturnCode should be 0")

		assert.Equal(t, WILLMSGRESP, msg.MessageType(), "MessageType() should return WILLMSGRESP")
	}
}
