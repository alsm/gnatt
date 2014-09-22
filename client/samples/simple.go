package main

import (
	"fmt"
	"github.com/alsm/gnatt/client"
	"github.com/alsm/gnatt/packets"
	"log"
	"os"
	"time"
)

func main() {
	gnatt.DEBUG = log.New(os.Stdout, "", 0)
	c, err := gnatt.NewClient("udp://127.0.0.1:1884", "gotestclient")
	if err != nil {
		panic(err)
	}
	c.SetWill("testwill", 1, true, []byte("Will test"))
	ct := c.Connect()
	ct.Wait()
	fmt.Println(ct.ReturnCode)
	rt := c.Register("gotest")
	rt.WaitTimeout(10 * time.Second)
	fmt.Println(rt.TopicName, rt.TopicId, rt.ReturnCode)
	st := c.Subscribe("gotest", 0, func(c *gnatt.SNClient, m *packets.PublishMessage) {
		fmt.Println(string(m.Data))
	})
	st.WaitTimeout(10 * time.Second)
	fmt.Println(st.TopicName, st.Qos, st.TopicId, st.ReturnCode)
	pt := c.Publish("gotest", 1, false, []byte("Hello gnatt"))
	pt.WaitTimeout(10 * time.Second)
	fmt.Println(pt.TopicId, pt.ReturnCode)
	time.Sleep(3 * time.Second)
}
