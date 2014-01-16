package main

type gwStatus byte

const (
	gw_stopped gwStatus = iota
	gw_starting
	gw_running
	gw_stopping
)

type gateway interface {
	start()
}
