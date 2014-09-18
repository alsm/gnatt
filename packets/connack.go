package packets

import (
	"io"
)

type ConnackMessage struct {
	Header
	ReturnCode byte
}

func (c *ConnackMessage) MessageType() byte {
	return CONNACK
}

func (c *ConnackMessage) Write(w io.Writer) error {
	packet := c.Header.pack()
	packet.WriteByte(CONNACK)
	packet.WriteByte(c.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (c *ConnackMessage) Unpack(b io.Reader) {
	c.ReturnCode = readByte(b)
}
