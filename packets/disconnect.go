package packets

import (
	"bytes"
	"io"
)

type DisconnectMessage struct {
	Header
	Duration uint16
}

func (d *DisconnectMessage) MessageType() byte {
	return DISCONNECT
}

func (d *DisconnectMessage) Write(w io.Writer) error {
	var packet bytes.Buffer

	if d.Duration == 0 {
		d.Header.Length = 2
		packet = d.Header.pack()
		packet.WriteByte(DISCONNECT)
	} else {
		d.Header.Length = 4
		packet = d.Header.pack()
		packet.WriteByte(DISCONNECT)
		packet.Write(encodeUint16(d.Duration))
	}
	_, err := packet.WriteTo(w)

	return err
}

func (d *DisconnectMessage) Unpack(b io.Reader) {
	if d.Header.Length == 4 {
		d.Duration = readUint16(b)
	}
}
