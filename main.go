package main

import (
	"fmt"
	"os"

	client "github.com/emil-nasso/share/client"
	server "github.com/emil-nasso/share/server"
)

//TODO: Flytta ut dessa konstanterna i communicator structen

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println(noCommandHelpMessage())
		os.Exit(1)
	}

	switch args[1] {
	case "upload":
		if len(args) != 3 {
			fmt.Println(missingGetArgumentHelpMessage())
			os.Exit(1)
		}
		client := client.New()
		client.Connect()
		sessionID := client.RequestUpload()
		fmt.Printf("SessionID: %v\n", sessionID)
		defer client.Disconnect()
		for {
			client.WaitAndSendFile(args[2])
		}
	case "server":
		server := server.New()
		server.Start()
	}
}
