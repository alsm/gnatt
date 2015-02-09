package packets

import (
	"io"
)

type SubscribeMessage struct {
	Header
	Dup         bool
	Qos         byte
	TopicIdType byte
	MessageId   uint16
	TopicId     uint16
	TopicName   []byte
}

func (s *SubscribeMessage) MessageType() byte {
	return SUBSCRIBE
}

func (s *SubscribeMessage) encodeFlags() byte {
	var b byte
	if s.Dup {
		b |= DUPFLAG
	}
	b |= (s.Qos << 5) & QOSBITS
	b |= s.TopicIdType & TOPICIDTYPE
	return b
}

func (s *SubscribeMessage) decodeFlags(b byte) {
	s.Dup = (b & DUPFLAG) == DUPFLAG
	s.Qos = (b & QOSBITS) >> 5
	s.TopicIdType = b & TOPICIDTYPE
}

func (s *SubscribeMessage) Write(w io.Writer) error {
	switch s.TopicIdType {
	case 0x00, 0x02:
		s.Header.Length = uint16(len(s.TopicName) + 5)
	case 0x01:
		s.Header.Length = 7
	}
	packet := s.Header.pack()
	packet.WriteByte(SUBSCRIBE)
	packet.WriteByte(s.encodeFlags())
	packet.Write(encodeUint16(s.MessageId))
	switch s.TopicIdType {
	case 0x00, 0x02:
		packet.Write(s.TopicName)
	case 0x01:
		packet.Write(encodeUint16(s.TopicId))
	}
	_, err := packet.WriteTo(w)

	return err
}

func (s *SubscribeMessage) Unpack(b io.Reader) {
	s.decodeFlags(readByte(b))
	s.MessageId = readUint16(b)
	switch s.TopicIdType {
	case 0x00, 0x02:
		s.TopicName = make([]byte, s.Header.Length-5)
		b.Read(s.TopicName)
	case 0x01:
		s.TopicId = readUint16(b)
	}
}
