package packets

import (
	"io"
)

type PubrecMessage struct {
	Header
	MessageId uint16
}

func (p *PubrecMessage) MessageType() byte {
	return PUBREC
}

func (p *PubrecMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PUBREC)
	packet.Write(encodeUint16(p.MessageId))
	_, err := packet.WriteTo(w)

	return err
}

func (p *PubrecMessage) Unpack(b io.Reader) {
	p.MessageId = readUint16(b)
}
