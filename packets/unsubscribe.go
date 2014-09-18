package packets

import (
	"io"
)

type UnsubscribeMessage struct {
	Header
	TopicIdType byte
	MessageId   uint16
	TopicId     uint16
	TopicName   []byte
}

func (u *UnsubscribeMessage) MessageType() byte {
	return UNSUBSCRIBE
}

func (s *UnsubscribeMessage) encodeFlags() byte {
	var b byte
	b |= s.TopicIdType & TOPICIDTYPE
	return b
}

func (s *UnsubscribeMessage) decodeFlags(b byte) {
	s.TopicIdType = b & TOPICIDTYPE
}

func (u *UnsubscribeMessage) Write(w io.Writer) error {
	switch u.TopicIdType {
	case 0x00, 0x02:
		u.Header.Length = uint16(len(u.TopicName) + 5)
	case 0x01:
		u.Header.Length = 7
	}
	packet := u.Header.pack()
	packet.WriteByte(UNSUBSCRIBE)
	packet.WriteByte(u.encodeFlags())
	packet.Write(encodeUint16(u.MessageId))
	switch u.TopicIdType {
	case 0x00, 0x02:
		packet.Write(u.TopicName)
	case 0x01:
		packet.Write(encodeUint16(u.TopicId))
	}
	_, err := packet.WriteTo(w)

	return err
}

func (u *UnsubscribeMessage) Unpack(b io.Reader) {
	u.decodeFlags(readByte(b))
	u.MessageId = readUint16(b)
	switch u.TopicIdType {
	case 0x00, 0x02:
		b.Read(u.TopicName)
	case 0x01:
		u.TopicId = readUint16(b)
	}
}
