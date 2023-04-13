package main

import (
	"github.com/XORbit01/DDOS-ARMY/cmd"
	"github.com/fatih/color"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			color.Red("Something went wrong: %v", r)
		}
	}()
	cmd.Execute()
}
