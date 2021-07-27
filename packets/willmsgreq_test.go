package packets

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWillMsgReqStruct(t *testing.T) {
	msg := NewMessage(WILLMSGREQ).(*WillMsgReqMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.WillMsgReqMessage", reflect.TypeOf(msg).String(), "Type should be WillMsgReqMessage")

		assert.Equal(t, WILLMSGREQ, msg.MessageType(), "MessageType() should return WILLMSGREQ")
	}
}
