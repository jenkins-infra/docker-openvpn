package cmd

import (
	"../easyrsa"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(requestCmd)
}

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Generate a client private key and a request certificate",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		errors := easyrsa.RequestClientCert(args)
		if errors != nil {
			for _, err := range errors {
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
	},
}
