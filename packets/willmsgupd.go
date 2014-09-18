package packets

import (
	"io"
)

type WillMsgUpdateMessage struct {
	Header
	WillMsg []byte
}

func (wm *WillMsgUpdateMessage) MessageType() byte {
	return WILLMSGUPD
}

func (wm *WillMsgUpdateMessage) Write(w io.Writer) error {
	wm.Header.Length = uint16(len(wm.WillMsg) + 2)
	packet := wm.Header.pack()
	packet.WriteByte(WILLMSGUPD)
	packet.Write(wm.WillMsg)
	_, err := packet.WriteTo(w)

	return err
}

func (wm *WillMsgUpdateMessage) Unpack(b io.Reader) {
	b.Read(wm.WillMsg)
}
