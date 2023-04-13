/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/XORbit01/DDOS-ARMY/client"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"strings"
)

// leaderCmd represents the leader command
var leaderCmd = &cobra.Command{
	Use:   "leader",
	Short: "connect as leader to server to control the army",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		connect, err := cmd.Flags().GetString("connect")
		if err != nil {
			color.Red("error getting host")
			return
		}
		if connect == "" {
			color.Red("please specify a server to connect using --connect flag or -c")
			return
		}
		if !strings.Contains(connect, ":") {
			connect = fmt.Sprintf("%s:8080", connect)
		}

		if !strings.HasPrefix(connect, "http://") {
			connect = fmt.Sprintf("http://%s", connect)
		}
		password, err := cmd.Flags().GetString("password")
		if err != nil {
			color.Red("error getting leader password, please check the password format")
			return
		}

		leader := client.Leader{
			Client:   client.Client{Name: "leader", DispatcherServer: connect},
			Password: password,
		}
		leader.Start()
	},
}

func init() {
	rootCmd.AddCommand(leaderCmd)
	leaderCmd.Flags().StringP("connect", "c", "", "host of the server")
	leaderCmd.Flags().StringP("password", "s", "", "secret password of the server")
}
