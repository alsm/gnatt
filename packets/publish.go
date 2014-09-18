package packets

import (
	"io"
)

type PublishMessage struct {
	Header
	Dup         bool
	Qos         byte
	TopicIdType byte
	Retain      bool
	TopicId     uint16
	MessageId   uint16
	Data        []byte
}

func (p *PublishMessage) MessageType() byte {
	return PUBLISH
}

func (p *PublishMessage) encodeFlags() byte {
	var b byte
	if p.Dup {
		b |= DUPFLAG
	}
	b |= (p.Qos << 5) & QOSBITS
	if p.Retain {
		b |= RETAINFLAG
	}
	b |= p.TopicIdType & TOPICIDTYPE
	return b
}

func (p *PublishMessage) decodeFlags(b byte) {
	p.Dup = (b & DUPFLAG) == DUPFLAG
	p.Qos = (b & QOSBITS) >> 5
	p.Retain = (b & RETAINFLAG) == RETAINFLAG
	p.TopicIdType = b & TOPICIDTYPE
}

func (p *PublishMessage) Write(w io.Writer) error {
	p.Header.Length = uint16(len(p.Data) + 7)
	packet := p.Header.pack()
	packet.WriteByte(PUBLISH)
	packet.WriteByte(p.encodeFlags())
	packet.Write(encodeUint16(p.TopicId))
	packet.Write(encodeUint16(p.MessageId))
	packet.Write(p.Data)
	_, err := packet.WriteTo(w)

	return err
}

func (p *PublishMessage) Unpack(b io.Reader) {
	p.decodeFlags(readByte(b))
	p.TopicId = readUint16(b)
	p.MessageId = readUint16(b)
	p.Data = make([]byte, p.Header.Length-7)
	b.Read(p.Data)
}
