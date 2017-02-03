package cmd

import (
	"github.com/emil-nasso/share/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve"},
	Short:   "Start the share server",
	Long: `The share server is a kind of connection broker. It waits for connections from uploaders,
creates a session and sends the uploader the session it. The session id can then be used by the downloader
to requests a download (either via the share binary or via http).
When this request comes, the server tells the uploader to start sending the file and simply forwards all bytes from the uploader to the downloader.
No files are stored on the server.`,
	Run: func(cmd *cobra.Command, args []string) {
		server := server.New()
		server.Start()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
