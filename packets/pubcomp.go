package packets

import (
	"io"
)

type PubcompMessage struct {
	Header
	MessageId uint16
}

func (p *PubcompMessage) MessageType() byte {
	return PUBCOMP
}

func (p *PubcompMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PUBCOMP)
	packet.Write(encodeUint16(p.MessageId))
	_, err := packet.WriteTo(w)

	return err
}

func (p *PubcompMessage) Unpack(b io.Reader) {
	p.MessageId = readUint16(b)
}
