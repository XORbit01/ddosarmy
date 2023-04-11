package cmd

import (
	"fmt"
	"github.com/XORbit01/DDOS-ARMY/client"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "join a ddos army camp!",
	Long: color.CyanString(`client is a command that joins a ddos army camp. 
you will be a soldier listening for commands from the dispatcher server.`),
	Run: func(cmd *cobra.Command, args []string) {
		//get the values from the flags
		connect, err := cmd.Flags().GetString("connect")
		if err != nil {
			color.Red("error getting connect flag")
			return
		}
		if !strings.Contains(connect, ":") {
			connect = fmt.Sprintf("%s:8080", connect)
		}

		if !strings.HasPrefix(connect, "http://") {
			connect = fmt.Sprintf("http://%s", connect)
		}
		if !IsValidAddr(connect) {
			//
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			color.Red("error getting name flag")
			return
		}

		cl := client.Client{Name: name, DispatcherServer: connect}
		cl.Start()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringP("connect", "c", "", "connect to a ddos army cam ")
	clientCmd.Flags().StringP("name", "n", getMachineHostName(), "your name, this will be used to identify you in the camp")
}
func getMachineHostName() string {
	name, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return name
}
