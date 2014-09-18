package packets

import (
	"io"
)

type PubackMessage struct {
	Header
	TopicId    uint16
	MessageId  uint16
	ReturnCode byte
}

func (p *PubackMessage) MessageType() byte {
	return PUBACK
}

func (p *PubackMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PUBACK)
	packet.Write(encodeUint16(p.TopicId))
	packet.Write(encodeUint16(p.MessageId))
	packet.WriteByte(p.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (p *PubackMessage) Unpack(b io.Reader) {
	p.TopicId = readUint16(b)
	p.MessageId = readUint16(b)
	p.ReturnCode = readByte(b)
}
