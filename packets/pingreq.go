package packets

import (
	"io"
)

type PingreqMessage struct {
	Header
	ClientId []byte
}

func (p *PingreqMessage) MessageType() byte {
	return PINGREQ
}

func (p *PingreqMessage) Write(w io.Writer) error {
	p.Header.Length = uint16(len(p.ClientId) + 2)
	packet := p.Header.pack()
	packet.WriteByte(PINGREQ)
	if len(p.ClientId) > 0 {
		packet.Write(p.ClientId)
	}
	_, err := packet.WriteTo(w)

	return err
}

func (p *PingreqMessage) Unpack(b io.Reader) {
	if p.Header.Length > 2 {
		b.Read(p.ClientId)
	}
}
