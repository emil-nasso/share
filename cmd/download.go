package cmd

import (
	"errors"
	"fmt"

	"github.com/emil-nasso/share/client"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <sessionid> <server hostname>",
	Short: "Download a file",
	Long: `The download command sends the <sessionid> to the server at <server hostname> to
request a download. The server tells the uploader to start the download and forwards the
bytes to the downloader.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			showErrorWithHelp(cmd, errors.New("Not enough arguments"))
		}

		sessionID := args[0]
		hostname := args[1]

		if sessionID == "" {
			showErrorWithHelp(cmd, errors.New("Invalid sessionID"))
		}
		fmt.Println("Downloading file for session", sessionID)
		client := client.New(hostname)
		defer client.Disconnect()
		client.RequestDownload(sessionID)
	},
}

func init() {
	RootCmd.AddCommand(downloadCmd)
}
