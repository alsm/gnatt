package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSubackStruct(t *testing.T) {
	msg := NewMessage(SUBACK).(*SubackMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.SubackMessage", reflect.TypeOf(msg).String(), "Type should be SubackMessage")
		assert.Equal(t, 8, msg.Length, "Default Length should be 8")
		assert.Equal(t, 0, msg.Qos, "Default Qos should be 0")
		assert.Equal(t, 0, msg.ReturnCode, "Default ReturnCode should be 0")
		assert.Equal(t, 0, msg.TopicId, "Default TopicId should be 0")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")

		assert.Equal(t, SUBACK, msg.MessageType(), "MessageType() should return SUBACK")
	}
}
