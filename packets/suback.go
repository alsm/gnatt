package packets

import (
	"io"
)

type SubackMessage struct {
	Header
	Qos        byte
	ReturnCode byte
	TopicId    uint16
	MessageId  uint16
}

func NewSubackMessage(TopicId uint16, MessageId uint16, Qos byte, rc byte) *SubackMessage {
	return &SubackMessage{
		Qos:        Qos,
		ReturnCode: rc,
		TopicId:    TopicId,
		MessageId:  MessageId,
	}
}

func (s *SubackMessage) MessageType() byte {
	return SUBACK
}

func (s *SubackMessage) encodeFlags() byte {
	var b byte
	b |= (s.Qos << 5) & QOSBITS
	return b
}

func (s *SubackMessage) decodeFlags(b byte) {
	s.Qos = (b & QOSBITS) >> 5
}

func (s *SubackMessage) Write(w io.Writer) error {
	packet := s.Header.pack()
	packet.WriteByte(SUBACK)
	packet.WriteByte(s.encodeFlags())
	packet.Write(encodeUint16(s.TopicId))
	packet.Write(encodeUint16(s.MessageId))
	packet.WriteByte(s.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (s *SubackMessage) Unpack(b io.Reader) {
	s.decodeFlags(readByte(b))
	s.TopicId = readUint16(b)
	s.MessageId = readUint16(b)
	s.ReturnCode = readByte(b)
}
