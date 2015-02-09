/*
Package gnatt/packets provides de/serialisation of MQTTSN packets
*/

package packets

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Message interface {
	MessageType() byte
	Write(io.Writer) error
	Unpack(io.Reader)
	//String() string
	//Details() Details
	//UUID() uuid.UUID
}

type Header struct {
	Length      uint16
	MessageType byte
}

func ReadPacket(r io.Reader) (m Message, err error) {
	var h Header
	packet := make([]byte, 1500)
	r.Read(packet)
	packetBuf := bytes.NewBuffer(packet)
	h.unpack(packetBuf)
	m = NewMessageWithHeader(h)
	if m == nil {
		return nil, errors.New("Bad data from client")
	}
	m.Unpack(packetBuf)
	return m, nil
}

func (h *Header) unpack(b io.Reader) {
	lengthCheck := readByte(b)
	if lengthCheck == 0x01 {
		h.Length = readUint16(b)
	} else {
		h.Length = uint16(lengthCheck)
	}
	h.MessageType = readByte(b)
}

func (h *Header) pack() bytes.Buffer {
	var header bytes.Buffer
	if h.Length > 256 {
		h.Length += 2
		header.WriteByte(0x01)
		header.Write(encodeUint16(h.Length))
	} else {
		header.WriteByte(byte(h.Length))
	}
	return header
}

func NewMessage(msgType byte) (m Message) {
	switch msgType {
	case ADVERTISE:
		m = &AdvertiseMessage{Header: Header{MessageType: ADVERTISE, Length: 5}}
	case SEARCHGW:
		m = &SearchGwMessage{Header: Header{MessageType: SEARCHGW, Length: 3}}
	case GWINFO:
		m = &GwInfoMessage{Header: Header{MessageType: GWINFO}}
	case CONNECT:
		m = &ConnectMessage{Header: Header{MessageType: CONNECT}, ProtocolId: 0x01}
	case CONNACK:
		m = &ConnackMessage{Header: Header{MessageType: CONNACK, Length: 3}}
	case WILLTOPICREQ:
		m = &WillTopicReqMessage{Header: Header{MessageType: WILLTOPICREQ, Length: 2}}
	case WILLTOPIC:
		m = &WillTopicMessage{Header: Header{MessageType: WILLTOPIC}}
	case WILLMSGREQ:
		m = &WillMsgReqMessage{Header: Header{MessageType: WILLMSGREQ, Length: 2}}
	case WILLMSG:
		m = &WillMsgMessage{Header: Header{MessageType: WILLMSG}}
	case REGISTER:
		m = &RegisterMessage{Header: Header{MessageType: REGISTER}}
	case REGACK:
		m = &RegackMessage{Header: Header{MessageType: REGACK, Length: 7}}
	case PUBLISH:
		m = &PublishMessage{Header: Header{MessageType: PUBLISH}}
	case PUBACK:
		m = &PubackMessage{Header: Header{MessageType: PUBACK, Length: 7}}
	case PUBCOMP:
		m = &PubcompMessage{Header: Header{MessageType: PUBCOMP, Length: 4}}
	case PUBREC:
		m = &PubrecMessage{Header: Header{MessageType: PUBREC, Length: 4}}
	case PUBREL:
		m = &PubrelMessage{Header: Header{MessageType: PUBREL, Length: 4}}
	case SUBSCRIBE:
		m = &SubscribeMessage{Header: Header{MessageType: SUBSCRIBE}}
	case SUBACK:
		m = &SubackMessage{Header: Header{MessageType: SUBACK, Length: 8}}
	case UNSUBSCRIBE:
		m = &UnsubscribeMessage{Header: Header{MessageType: UNSUBSCRIBE}}
	case UNSUBACK:
		m = &UnsubackMessage{Header: Header{MessageType: UNSUBACK, Length: 4}}
	case PINGREQ:
		m = &PingreqMessage{Header: Header{MessageType: PINGREQ}}
	case PINGRESP:
		m = &PingrespMessage{Header: Header{MessageType: PINGRESP, Length: 2}}
	case DISCONNECT:
		m = &DisconnectMessage{Header: Header{MessageType: DISCONNECT}}
	case WILLTOPICUPD:
		m = &WillTopicUpdateMessage{Header: Header{MessageType: WILLTOPICUPD}}
	case WILLTOPICRESP:
		m = &WillTopicRespMessage{Header: Header{MessageType: WILLTOPICRESP, Length: 3}}
	case WILLMSGUPD:
		m = &WillMsgUpdateMessage{Header: Header{MessageType: WILLMSGUPD}}
	case WILLMSGRESP:
		m = &WillMsgRespMessage{Header: Header{MessageType: WILLMSGRESP, Length: 3}}
	}
	return
}

func NewMessageWithHeader(h Header) (m Message) {
	switch h.MessageType {
	case ADVERTISE:
		m = &AdvertiseMessage{Header: h}
	case SEARCHGW:
		m = &SearchGwMessage{Header: h}
	case GWINFO:
		m = &GwInfoMessage{Header: h}
	case CONNECT:
		m = &ConnectMessage{Header: h}
	case CONNACK:
		m = &ConnackMessage{Header: h}
	case WILLTOPICREQ:
		m = &WillTopicReqMessage{Header: h}
	case WILLTOPIC:
		m = &WillTopicMessage{Header: h}
	case WILLMSGREQ:
		m = &WillMsgReqMessage{Header: h}
	case WILLMSG:
		m = &WillMsgMessage{Header: h}
	case REGISTER:
		m = &RegisterMessage{Header: h}
	case REGACK:
		m = &RegackMessage{Header: h}
	case PUBLISH:
		m = &PublishMessage{Header: h}
	case PUBACK:
		m = &PubackMessage{Header: h}
	case PUBCOMP:
		m = &PubcompMessage{Header: h}
	case PUBREC:
		m = &PubrecMessage{Header: h}
	case PUBREL:
		m = &PubrelMessage{Header: h}
	case SUBSCRIBE:
		m = &SubscribeMessage{Header: h}
	case SUBACK:
		m = &SubackMessage{Header: h}
	case UNSUBSCRIBE:
		m = &UnsubscribeMessage{Header: h}
	case UNSUBACK:
		m = &UnsubackMessage{Header: h}
	case PINGREQ:
		m = &PingreqMessage{Header: h}
	case PINGRESP:
		m = &PingrespMessage{Header: h}
	case DISCONNECT:
		m = &DisconnectMessage{Header: h}
	case WILLTOPICUPD:
		m = &WillTopicUpdateMessage{Header: h}
	case WILLTOPICRESP:
		m = &WillTopicRespMessage{Header: h}
	case WILLMSGUPD:
		m = &WillMsgUpdateMessage{Header: h}
	case WILLMSGRESP:
		m = &WillMsgRespMessage{Header: h}
	}
	return
}

func readByte(b io.Reader) byte {
	num := make([]byte, 1)
	b.Read(num)
	return num[0]
}

func readUint16(b io.Reader) uint16 {
	num := make([]byte, 2)
	b.Read(num)
	return binary.BigEndian.Uint16(num)
}

func encodeUint16(num uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, num)
	return bytes
}

// Flags
const (
	TOPICIDTYPE  = 0x03
	CLEANSESSION = 0x04
	WILLFLAG     = 0x08
	RETAINFLAG   = 0x10
	QOSBITS      = 0x60
	DUPFLAG      = 0x80
)

// Errors
const (
	ACCEPTED         = 0x00
	REJ_CONGESTION   = 0x01
	REJ_INVALID_TID  = 0x02
	REJ_NOT_SUPORTED = 0x03
)

// Message Types
const (
	ADVERTISE     = 0x00
	SEARCHGW      = 0x01
	GWINFO        = 0x02
	CONNECT       = 0x04
	CONNACK       = 0x05
	WILLTOPICREQ  = 0x06
	WILLTOPIC     = 0x07
	WILLMSGREQ    = 0x08
	WILLMSG       = 0x09
	REGISTER      = 0x0A
	REGACK        = 0x0B
	PUBLISH       = 0x0C
	PUBACK        = 0x0D
	PUBCOMP       = 0x0E
	PUBREC        = 0x0F
	PUBREL        = 0x10
	SUBSCRIBE     = 0x12
	SUBACK        = 0x13
	UNSUBSCRIBE   = 0x14
	UNSUBACK      = 0x15
	PINGREQ       = 0x16
	PINGRESP      = 0x17
	DISCONNECT    = 0x18
	WILLTOPICUPD  = 0x1A
	WILLTOPICRESP = 0x1B
	WILLMSGUPD    = 0x1C
	WILLMSGRESP   = 0x1D
	// 0x03 is reserved
	// 0x11 is reserved
	// 0x19 is reserved
	// 0x1E - 0xFD is reserved
	// 0xFE - Encapsulated message
	// 0xFF is reserved
)

var MessageNames = map[byte]string{
	ADVERTISE:     "ADVERTISE",
	SEARCHGW:      "SEARCHGW",
	GWINFO:        "GWINFO",
	CONNECT:       "CONNECT",
	CONNACK:       "CONNACK",
	WILLTOPICREQ:  "WILLTOPICREQ",
	WILLTOPIC:     "WILLTOPIC",
	WILLMSGREQ:    "WILLMSGREQ",
	WILLMSG:       "WILLMSG",
	REGISTER:      "REGISTER",
	REGACK:        "REGACK",
	PUBLISH:       "PUBLISH",
	PUBACK:        "PUBACK",
	PUBCOMP:       "PUBCOMP",
	PUBREC:        "PUBREC",
	PUBREL:        "PUBREL",
	SUBSCRIBE:     "SUBSCRIBE",
	SUBACK:        "SUBACK",
	UNSUBSCRIBE:   "UNSUBSCRIBE",
	UNSUBACK:      "UNSUBACK",
	PINGREQ:       "PINGREQ",
	PINGRESP:      "PINGRESP",
	DISCONNECT:    "DISCONNECT",
	WILLTOPICUPD:  "WILLTOPICUPD",
	WILLTOPICRESP: "WILLTOPICRESP",
	WILLMSGUPD:    "WILLMSGUPD",
	WILLMSGRESP:   "WILLMSGRESP",
}
