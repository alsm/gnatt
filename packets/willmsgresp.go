package packets

import (
	"io"
)

type WillMsgRespMessage struct {
	Header
	ReturnCode byte
}

func (wm *WillMsgRespMessage) MessageType() byte {
	return WILLMSGRESP
}

func (wm *WillMsgRespMessage) Write(w io.Writer) error {
	packet := wm.Header.pack()
	packet.WriteByte(WILLMSGRESP)
	packet.WriteByte(wm.ReturnCode)
	_, err := packet.WriteTo(w)

	return err
}

func (wm *WillMsgRespMessage) Unpack(b io.Reader) {
	wm.ReturnCode = readByte(b)
}
