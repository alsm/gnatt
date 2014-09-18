package packets

import (
	"io"
)

type WillMsgReqMessage struct {
	Header
}

func (wm *WillMsgReqMessage) MessageType() byte {
	return WILLMSGREQ
}

func (wm *WillMsgReqMessage) Write(w io.Writer) error {
	packet := wm.Header.pack()
	packet.WriteByte(wm.Header.MessageType)
	_, err := packet.WriteTo(w)

	return err
}

func (wm *WillMsgReqMessage) Unpack(b io.Reader) {

}
