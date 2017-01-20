package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	client "github.com/emil-nasso/share/client"
	server "github.com/emil-nasso/share/server"
)

//TODO: Flytta ut dessa konstanterna i communicator structen

func main() {
	rand.Seed(time.Now().UnixNano())

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
		fmt.Printf("SessionID: %v\n\n", sessionID)
		fmt.Printf("Download url:\n http://%s:%d/get/%s\n", client.ServerHostname, 27002, sessionID)
		fmt.Printf("Download with share:\n share download %v\n", sessionID)
		defer client.Disconnect()
		for {
			client.WaitAndSendFile(args[2])
		}
	case "download":
		if len(args) == 3 {
			var sessionID string
			sessionID = args[2]
			fmt.Println("Downloading file for session", sessionID)
			client := client.New()
			client.Connect()
			client.RequestDownload(sessionID)
		}
	case "server":
		server := server.New()
		server.Start()
	default:
		fmt.Println("invalid usage")
	}
}
