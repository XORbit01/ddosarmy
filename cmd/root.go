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
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
