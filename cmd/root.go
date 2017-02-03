package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/emil-nasso/share/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "share",
	Short: "An application for sharing files, even if you are begi",
	Long: `Share consists of three components, the server, uploader and downloader.
More information about each of these parts are available in each commands help page.
You can use the help command to read about other commands.

Ex: share help server`,
}

// Execute is called from main and initiates the commands
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func showErrorWithHelp(cmd *cobra.Command, err error) {
	fmt.Println(err)
	cmd.Help()
	os.Exit(1)
}

func onConfigLoaded() {
	lib.DebugEnabled = viper.GetBool("debug")
	lib.Debug("Debug is enabled")
}

func init() {
	rand.Seed(time.Now().UnixNano())
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test-cobra.yaml)")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Show debug information")
	viper.BindPFlag("debug", RootCmd.PersistentFlags().Lookup("debug"))
}

func initConfig() {
	viper.SetConfigName(".test-cobra")
	viper.AddConfigPath("$HOME")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println(err)
	}
	onConfigLoaded()
}
