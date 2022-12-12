package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "easyvpn",
	Short: "Easyvpn is a client tool to manage Jenkins Infrastructure VPN",
}

// CertDir define cert directory path
var CertDir string

// Net define network to use (determine the ccd to use)
var Network string

// Commit define if changes must be committed or not
var Commit bool

// Push define if changes must be pushed or not
var Push bool

// Execute run the main command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
