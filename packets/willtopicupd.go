package packets

import (
	"io"
)

type WillTopicUpdateMessage struct {
	Header
	Qos       byte
	Retain    bool
	WillTopic []byte
}

func (wt *WillTopicUpdateMessage) MessageType() byte {
	return WILLTOPICUPD
}

func (wt *WillTopicUpdateMessage) encodeFlags() byte {
	var b byte
	b |= (wt.Qos << 5) & QOSBITS
	if wt.Retain {
		b |= RETAINFLAG
	}
	return b
}

func (wt *WillTopicUpdateMessage) decodeFlags(b byte) {
	wt.Qos = (b & QOSBITS) >> 5
	wt.Retain = (b & RETAINFLAG) == RETAINFLAG
}

func (wt *WillTopicUpdateMessage) Write(w io.Writer) error {
	wt.Header.Length = uint16(len(wt.WillTopic) + 3)
	packet := wt.Header.pack()
	packet.WriteByte(WILLTOPICUPD)
	packet.WriteByte(wt.encodeFlags())
	packet.Write(wt.WillTopic)
	_, err := packet.WriteTo(w)

	return err
}

func (wt *WillTopicUpdateMessage) Unpack(b io.Reader) {
	wt.decodeFlags(readByte(b))
	b.Read(wt.WillTopic)
}
