package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestPubackStruct(t *testing.T) {
	msg := NewMessage(PUBACK).(*PubackMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.PubackMessage", reflect.TypeOf(msg).String(), "Type should be PubackMessage")
		assert.Equal(t, 0, msg.TopicId, "Default TopicId should be 0")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")
		assert.Equal(t, 0, msg.ReturnCode, "Default ReturnCode should be 0")
		assert.Equal(t, 7, msg.Length, "Default Length should be 2")

		assert.Equal(t, PUBACK, msg.MessageType(), "MessageType() should return PUBACK")
	}
}
