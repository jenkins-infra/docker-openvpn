package cmd

import (
	"../easyrsa"
	"../git"
	"../helpers"
	"fmt"
	"github.com/spf13/cobra"
	"path"
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().BoolVarP(&Commit, "commit", "", true, "git commit changes")
	signCmd.Flags().BoolVarP(&Push, "push", "", true, "git push changes")
	signCmd.Flags().StringVarP(&CertDir, "cert", "c", "cert", "Cert Directory")
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

		if Commit {
			for i := 0; i < len(args); i++ {
				msg := "[infra-admin] Sign certificate request for " + args[i]
				files := []string{
					path.Join(CertDir, "pki/issued", args[i]+".crt"),
					path.Join(CertDir, "pki", "index.txt"),
					path.Join(CertDir, "pki", "serial"),
				}
				git.Add(files)
				git.Commit(files, msg)
			}
		}

		if Push {
			git.Push()
		}

		helpers.CleanPrivateDir()
	},
}
