package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "easyvpn",
	Short: "Easyvpn is a client tool to manage Jenkins Infrastructure VPN",
}

// Execute run the main command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
