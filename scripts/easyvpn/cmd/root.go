package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "easyvpn",
	Short: "Easyvpn is a client tool to manage Jenkins Infrastructure VPN",
}

// Execute run the main command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
