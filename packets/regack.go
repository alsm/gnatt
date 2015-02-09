package packets

import (
	"io"
)

type RegackMessage struct {
	Header
	TopicId    uint16
	MessageId  uint16
	ReturnCode byte
}

func NewRegackMessage(TopicId uint16, MessageId uint16, rc byte) *RegackMessage {
	return &RegackMessage{
		TopicId:    TopicId,
		MessageId:  MessageId,
		ReturnCode: rc,
	}
}

func (r *RegackMessage) MessageType() byte {
	return REGACK
}

func (r *RegackMessage) Write(w io.Writer) error {
	packet := r.Header.pack()
	packet.WriteByte(REGACK)
	packet.Write(encodeUint16(r.TopicId))
	packet.Write(encodeUint16(r.MessageId))
	packet.WriteByte(r.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (r *RegackMessage) Unpack(b io.Reader) {
	r.TopicId = readUint16(b)
	r.MessageId = readUint16(b)
	r.ReturnCode = readByte(b)
}
