package packets

import (
	"io"
)

type RegisterMessage struct {
	Header
	TopicId   uint16
	MessageId uint16
	TopicName []byte
}

func (r *RegisterMessage) MessageType() byte {
	return REGISTER
}

func (r *RegisterMessage) Write(w io.Writer) error {
	r.Header.Length = uint16(len(r.TopicName) + 6)
	packet := r.Header.pack()
	packet.WriteByte(REGISTER)
	packet.Write(encodeUint16(r.TopicId))
	packet.Write(encodeUint16(r.MessageId))
	packet.Write(r.TopicName)
	_, err := packet.WriteTo(w)

	return err
}

func (r *RegisterMessage) Unpack(b io.Reader) {
	r.TopicId = readUint16(b)
	r.MessageId = readUint16(b)
	b.Read(r.TopicName)
}
