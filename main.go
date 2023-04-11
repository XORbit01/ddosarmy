package main

import (
	"github.com/XORbit01/DDOS-ARMY/cmd"
	"log"
	"os"
)

func main() {
	file, err := os.Create("logfile.txt")
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}
	defer file.Close()

	// Set the output of the logger to the file
	log.SetOutput(file)

	cmd.Execute()
}
