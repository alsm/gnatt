package packets

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessageTypeValues(t *testing.T) {
	assert.Equal(t, 0x00, ADVERTISE, "ADVERTISE should be 0x00")
	assert.Equal(t, 0x01, SEARCHGW, "SEARCHGW should be 0x01")
	assert.Equal(t, 0x02, GWINFO, "GWINFO should be 0x02")
	assert.Equal(t, 0x04, CONNECT, "CONNECT should be 0x04")
	assert.Equal(t, 0x05, CONNACK, "CONNACK should be 0x05")
	assert.Equal(t, 0x06, WILLTOPICREQ, "WILLTOPICREQ should be 0x06")
	assert.Equal(t, 0x07, WILLTOPIC, "WILLTOPIC should be 0x07")
	assert.Equal(t, 0x08, WILLMSGREQ, "WILLMSGREQ should be 0x08")
	assert.Equal(t, 0x09, WILLMSG, "WILLMSG should be 0x09")
	assert.Equal(t, 0x0A, REGISTER, "REGISTER should be 0x0A")
	assert.Equal(t, 0x0B, REGACK, "REGACK should be 0x0B")
	assert.Equal(t, 0x0C, PUBLISH, "PUBLISH should be 0x0C")
	assert.Equal(t, 0x0D, PUBACK, "PUBACK should be 0x0D")
	assert.Equal(t, 0x0E, PUBCOMP, "PUBCOMP should be 0x0E")
	assert.Equal(t, 0x0F, PUBREC, "PUBREC should be 0x0F")
	assert.Equal(t, 0x10, PUBREL, "PUBREL should be 0x10")
	assert.Equal(t, 0x12, SUBSCRIBE, "SUBSCRIBE should be 0x12")
	assert.Equal(t, 0x13, SUBACK, "SUBACK should be 0x13")
	assert.Equal(t, 0x14, UNSUBSCRIBE, "UNSUBSCRIBE should be 0x14")
	assert.Equal(t, 0x15, UNSUBACK, "UNSUBACK should be 0x15")
	assert.Equal(t, 0x16, PINGREQ, "PINGREQ should be 0x16")
	assert.Equal(t, 0x17, PINGRESP, "PINGRESPshould be 0x17")
	assert.Equal(t, 0x18, DISCONNECT, "DISCONNECTshould be 0x18")
	assert.Equal(t, 0x1A, WILLTOPICUPD, "WILLTOPICUPDshould be 0x1A")
	assert.Equal(t, 0x1B, WILLTOPICRESP, "WILLTOPICRESshould be 0x1B")
	assert.Equal(t, 0x1C, WILLMSGUPD, "WILLMSGUPDshould be 0x1C")
	assert.Equal(t, 0x1D, WILLMSGRESP, "WILLMSGRESPshould be 0x1D")
}
