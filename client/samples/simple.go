package main

import (
	"fmt"
	"github.com/alsm/gnatt/client"
	"log"
	"os"
	"time"
)

func main() {
	gnatt.DEBUG = log.New(os.Stdout, "", 0)
	c := gnatt.NewClient("udp://127.0.0.1:1884")
	c.Connect()
	time.Sleep(1 * time.Second)
	rt := c.Register("gotest")
	rt.WaitTimeout(10 * time.Second)
	fmt.Println(rt.TopicName, rt.TopicId, rt.ReturnCode)
	st := c.Subscribe("gotest", 0)
	st.WaitTimeout(10 * time.Second)
	fmt.Println(st.TopicName, st.Qos, st.TopicId, st.ReturnCode)
}
