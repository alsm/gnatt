package gateway

type Gateway interface {
	Start()
	Port() int
	OnPacket(int, []byte, uConn, uAddr)
}
