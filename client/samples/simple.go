package main

import (
	"github.com/alsm/gnatt/client"
	"log"
	"os"
	"time"
)

func main() {
	gnatt.DEBUG = log.New(os.Stdout, "", 0)
	c := gnatt.NewClient("udp://127.0.0.1:1884")
	c.Connect()
	time.Sleep(5 * time.Second)
}
