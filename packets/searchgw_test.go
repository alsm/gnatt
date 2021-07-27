package packets

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchGwStruct(t *testing.T) {
	msg := NewMessage(SEARCHGW).(*SearchGwMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.SearchGwMessage", reflect.TypeOf(msg).String(), "Type should be SearchGwMessage")
		assert.Equal(t, 3, msg.Length, "Default Length should be 3")
		assert.Equal(t, 0, msg.Radius, "Default Radius should be 0")

		assert.Equal(t, SEARCHGW, msg.MessageType(), "MessageType() should return SEARCHGW")
	}
}
