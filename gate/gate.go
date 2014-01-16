package gateway

type gwStatus byte

const (
	gw_stopped gwStatus = iota
	gw_starting
	gw_running
	gw_stopping
)

type Gateway interface {
	Start()
}
