package packets

import (
	"io"
)

type ConnectMessage struct {
	Header
	Will         bool
	CleanSession bool
	ProtocolId   byte
	Duration     uint16
	ClientId     []byte
}

func (c *ConnectMessage) MessageType() byte {
	return CONNECT
}

func (c *ConnectMessage) decodeFlags(b byte) {
	c.Will = (b & WILLFLAG) == WILLFLAG
	c.CleanSession = (b & CLEANSESSION) == CLEANSESSION
}

func (c *ConnectMessage) encodeFlags() byte {
	var b byte
	if c.Will {
		b |= WILLFLAG
	}
	if c.CleanSession {
		b |= CLEANSESSION
	}
	return b
}

func (c *ConnectMessage) Write(w io.Writer) error {
	c.Header.Length = uint16(len(c.ClientId) + 6)
	packet := c.Header.pack()
	packet.WriteByte(CONNECT)
	packet.WriteByte(c.encodeFlags())
	packet.WriteByte(c.ProtocolId)
	packet.Write(encodeUint16(c.Duration))
	packet.Write([]byte(c.ClientId))
	_, err := packet.WriteTo(w)

	return err
}

func (c *ConnectMessage) Unpack(b io.Reader) {
	c.decodeFlags(readByte(b))
	c.ProtocolId = readByte(b)
	c.Duration = readUint16(b)
	c.ClientId = make([]byte, c.Header.Length-6)
	b.Read(c.ClientId)
}
