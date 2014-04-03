package gateway

import (
	"io"
	"log"
)

var (
	INFO  *log.Logger
	ERROR *log.Logger
)

func InitLogger(infoHandle, errorHandle io.Writer) {
	INFO = log.New(infoHandle, "INFO:  ", log.Ldate|log.Ltime)
	ERROR = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime)
}
