package packets

import (
	"io"
)

type GwInfoMessage struct {
	Header
	GatewayId      byte
	GatewayAddress []byte
}

func (g *GwInfoMessage) MessageType() byte {
	return GWINFO
}

func (g *GwInfoMessage) Write(w io.Writer) error {
	g.Header.Length = uint16(len(g.GatewayAddress) + 3)
	packet := g.Header.pack()
	packet.WriteByte(GWINFO)
	packet.WriteByte(g.GatewayId)
	packet.Write(g.GatewayAddress)
	_, err := packet.WriteTo(w)

	return err
}

func (g *GwInfoMessage) Unpack(b io.Reader) {
	g.GatewayId = readByte(b)
	if g.Header.Length > 3 {
		b.Read(g.GatewayAddress)
	}
}
