package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestWillTopicStruct(t *testing.T) {
	msg := NewMessage(WILLTOPIC).(*WillTopicMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.WillTopicMessage", reflect.TypeOf(msg).String(), "Type should be WillTopicMessage")
		assert.Equal(t, 0, msg.Qos, "Default Qos should be 0")
		assert.Equal(t, false, msg.Retain, "Default Retain flag should be false")
		assert.Equal(t, []byte(nil), msg.WillTopic, "Default WillMsg should be blank")

		assert.Equal(t, WILLTOPIC, msg.MessageType(), "MessageType() should return WILLTOPIC")
	}
}
