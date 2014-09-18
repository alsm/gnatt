package packets

import (
	"io"
)

type PingrespMessage struct {
	Header
}

func (p *PingrespMessage) MessageType() byte {
	return PINGRESP
}

func (p *PingrespMessage) Write(w io.Writer) error {
	packet := p.Header.pack()
	packet.WriteByte(PINGRESP)
	_, err := packet.WriteTo(w)

	return err
}

func (p *PingrespMessage) Unpack(b io.Reader) {
}
