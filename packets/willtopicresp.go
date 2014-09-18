package packets

import (
	"io"
)

type WillTopicRespMessage struct {
	Header
	ReturnCode byte
}

func (wt *WillTopicRespMessage) MessageType() byte {
	return WILLTOPICRESP
}

func (wt *WillTopicRespMessage) Write(w io.Writer) error {
	packet := wt.Header.pack()
	packet.WriteByte(WILLTOPICRESP)
	packet.WriteByte(wt.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (wt *WillTopicRespMessage) Unpack(b io.Reader) {
	wt.ReturnCode = readByte(b)
}
