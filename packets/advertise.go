package packets

import (
	"io"
)

type AdvertiseMessage struct {
	Header
	GatewayId byte
	Duration  uint16
}

func (a *AdvertiseMessage) MessageType() byte {
	return ADVERTISE
}

func (a *AdvertiseMessage) Write(w io.Writer) error {
	packet := a.Header.pack()
	packet.WriteByte(ADVERTISE)
	packet.WriteByte(a.GatewayId)
	packet.Write(encodeUint16(a.Duration))
	_, err := packet.WriteTo(w)

	return err
}

func (a *AdvertiseMessage) Unpack(b io.Reader) {
	a.GatewayId = readByte(b)
	a.Duration = readUint16(b)
}
