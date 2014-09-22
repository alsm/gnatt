package packets

import (
	"io"
)

type WillTopicMessage struct {
	Header
	Qos       byte
	Retain    bool
	WillTopic []byte
}

func (wt *WillTopicMessage) MessageType() byte {
	return wt.Header.MessageType
}

func (wt *WillTopicMessage) encodeFlags() byte {
	var b byte

	b |= (wt.Qos << 5) & QOSBITS
	if wt.Retain {
		b |= RETAINFLAG
	}
	return b
}

func (wt *WillTopicMessage) decodeFlags(b byte) {
	wt.Qos = (b & QOSBITS) >> 5
	wt.Retain = (b & RETAINFLAG) == RETAINFLAG
}

func (wt *WillTopicMessage) Write(w io.Writer) error {
	if len(wt.WillTopic) == 0 {
		wt.Header.Length = 2
	} else {
		wt.Header.Length = uint16(len(wt.WillTopic) + 3)
	}
	packet := wt.Header.pack()
	packet.WriteByte(wt.Header.MessageType)
	if wt.Header.Length > 2 {
		packet.WriteByte(wt.encodeFlags())
		packet.Write(wt.WillTopic)
	}
	_, err := packet.WriteTo(w)

	return err
}

func (wt *WillTopicMessage) Unpack(b io.Reader) {
	if wt.Header.Length > 2 {
		wt.decodeFlags(readByte(b))
		b.Read(wt.WillTopic)
	}
}
