package cmd

import (
	"../easyrsa"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(revokeCmd)
}

var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke a client certificate",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		errors := easyrsa.RevokeClientCert(args)
		if errors != nil {
			for _, err := range errors {
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
	},
}
