package packets

import (
	"io"
)

type WillTopicReqMessage struct {
	Header
}

func (wt *WillTopicReqMessage) MessageType() byte {
	return WILLTOPICREQ
}

func (wt *WillTopicReqMessage) Write(w io.Writer) error {
	packet := wt.Header.pack()
	packet.WriteByte(WILLTOPICREQ)
	_, err := packet.WriteTo(w)

	return err
}

func (wt *WillTopicReqMessage) Unpack(b io.Reader) {

}
