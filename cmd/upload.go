package cmd

import (
	"fmt"
	"os"

	"github.com/emil-nasso/share/client"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:        "upload <filepath> <hostname>",
	Aliases:    []string{"up", "send"},
	ArgAliases: []string{"filepath", "hostname"},
	Short:      "Initiate an upload request and wait for downloaders.",
	Long: `When running the upload command share connects to the server (specified by hostname) and requests a session id.
This session id can then be sent to the downloader and can be used to download the file (specified by filepath)`,
	ValidArgs: []string{"file", "server"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Help()
			os.Exit(1)
		}

		filePath := args[0]
		serverHostname := args[1]

		if filePath == "" {
			fmt.Println("Could not open file:", filePath)
			cmd.Help()
			os.Exit(1)
		}
		client := client.New(serverHostname)

		sessionID := client.RequestUpload()
		fmt.Printf("SessionID: %v\n", sessionID)
		fmt.Printf("Download url:\n http://%s:%d/get/%s\n", client.ServerHostname, 27002, sessionID)
		fmt.Printf("Download with share:\n share download %v %v\n", sessionID, client.ServerHostname)
		defer client.Disconnect()
		for {
			client.WaitAndSendFile(filePath)
		}
	},
}

func init() {
	RootCmd.AddCommand(uploadCmd)
}
