package packets

import (
	"io"
)

type UnsubackMessage struct {
	Header
	MessageId uint16
}

func (u *UnsubackMessage) MessageType() byte {
	return UNSUBACK
}

func (u *UnsubackMessage) Write(w io.Writer) error {
	packet := u.Header.pack()
	packet.WriteByte(UNSUBACK)
	packet.Write(encodeUint16(u.MessageId))
	_, err := packet.WriteTo(w)

	return err
}

func (u *UnsubackMessage) Unpack(b io.Reader) {
	u.MessageId = readUint16(b)
}
