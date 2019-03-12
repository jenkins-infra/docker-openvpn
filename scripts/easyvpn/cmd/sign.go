package cmd

import (
	"../easyrsa"
	"../helpers"
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(signCmd)
}

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a request certificate",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		helpers.DecryptPrivateDir()
		errors := easyrsa.SignClientRequest(args)
		if errors != nil {
			for _, err := range errors {
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
		helpers.CleanPrivateDir()
	},
}
