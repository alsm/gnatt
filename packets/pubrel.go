package packets

import (
	"io"
)

type PubrelMessage struct {
	Header
	MessageId uint16
}

func (p *PubrelMessage) MessageType() byte {
	return PUBREL
}

func (p *PubrelMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PUBREL)
	packet.Write(encodeUint16(p.MessageId))
	_, err := packet.WriteTo(w)

	return err
}

func (p *PubrelMessage) Unpack(b io.Reader) {
	p.MessageId = readUint16(b)
}
