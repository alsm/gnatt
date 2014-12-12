package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestWillTopicUpdateStruct(t *testing.T) {
	msg := NewMessage(WILLTOPICUPD).(*WillTopicUpdateMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.WillTopicUpdateMessage", reflect.TypeOf(msg).String(), "Type should be WillTopicUpdateMessage")
		assert.Equal(t, 0, msg.Qos, "Default Qos should be 0")
		assert.Equal(t, false, msg.Retain, "Default Retain flag should be false")
		assert.Equal(t, []byte(nil), msg.WillTopic, "Default WillMsg should be blank")

		assert.Equal(t, WILLTOPICUPD, msg.MessageType(), "MessageType() should return WILLTOPICUPD")
	}

}
