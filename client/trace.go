package gnatt

import (
	"io/ioutil"
	"log"
)

var (
	ERROR    *log.Logger
	CRITICAL *log.Logger
	WARN     *log.Logger
	DEBUG    *log.Logger
)

const (
	NET = "[net]     "
	PNG = "[pinger]  "
	CLI = "[client]  "
	DEC = "[decode]  "
	MES = "[message] "
	STR = "[store]   "
	MID = "[msgids]  "
	TST = "[test]    "
	STA = "[state]   "
	ERR = "[error]   "
)

func init() {
	ERROR = log.New(ioutil.Discard, "", 0)
	CRITICAL = log.New(ioutil.Discard, "", 0)
	WARN = log.New(ioutil.Discard, "", 0)
	DEBUG = log.New(ioutil.Discard, "", 0)
}
