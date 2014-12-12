package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUnsubackStruct(t *testing.T) {
	msg := NewMessage(UNSUBACK).(*UnsubackMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.UnsubackMessage", reflect.TypeOf(msg).String(), "Type should be UnsubackMessage")
		assert.Equal(t, 4, msg.Length, "Default Length should be 4")
		assert.Equal(t, 0, msg.MessageId, "Default MessageId should be 0")

		assert.Equal(t, UNSUBACK, msg.MessageType(), "MessageType() should return UNSUBACK")
	}
}
