package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	client "github.com/emil-nasso/share/client"
	"github.com/emil-nasso/share/lib"
	server "github.com/emil-nasso/share/server"
)

//TODO: Flytta ut dessa konstanterna i communicator structen

func main() {
	rand.Seed(time.Now().UnixNano())

	hostname := flag.String("s", "localhost", "the hostname/ip of the server to connect to")
	lib.DebugEnabled = *flag.Bool("d", false, "show debug information")
	flag.Parse()
	command := flag.Arg(0)

	if command == "" {
		fmt.Println(noCommandHelpMessage())
		os.Exit(1)
	}

	switch command {
	case "upload":
		filePath := flag.Arg(1)
		if filePath == "" {
			fmt.Println(missingGetArgumentHelpMessage())
			os.Exit(1)
		}
		client := client.New(*hostname)
		sessionID := client.RequestUpload()
		fmt.Printf("SessionID: %v\n", sessionID)
		fmt.Printf("Download url:\n http://%s:%d/get/%s\n", client.ServerHostname, 27002, sessionID)
		fmt.Printf("Download with share:\n share download %v\n", sessionID)
		defer client.Disconnect()
		for {
			client.WaitAndSendFile(filePath)
		}
	case "download":
		sessionID := flag.Arg(1)
		if sessionID == "" {
			fmt.Println("No sessionID")
			os.Exit(1)
		}
		fmt.Println("Downloading file for session", sessionID)
		client := client.New(*hostname)
		defer client.Disconnect()
		client.RequestDownload(sessionID)
	case "server":
		server := server.New()
		server.Start()
	default:
		fmt.Println("invalid usage")
	}
}
