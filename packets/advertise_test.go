package packets

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestAdvertiseStruct(t *testing.T) {
	msg := NewMessage(ADVERTISE).(*AdvertiseMessage)

	if assert.NotNil(t, msg, "New message should not be nil") {
		assert.Equal(t, "*packets.AdvertiseMessage", reflect.TypeOf(msg).String(), "Type should be AdvertiseMessage")
		assert.Equal(t, 0, msg.GatewayId, "Default GatewayId should be 0")
		assert.Equal(t, 0, msg.Duration, "Default Duration should be 0")
		assert.Equal(t, 5, msg.Length, "Length should be 5")

		assert.Equal(t, ADVERTISE, msg.MessageType(), "MessageType() should return ADVERTISE")
	}
}
