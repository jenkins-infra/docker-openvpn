package cmd

import (
	"../easyrsa"
	"../git"
	"fmt"
	"github.com/spf13/cobra"
	"path"
)

func init() {
	rootCmd.AddCommand(revokeCmd)
	revokeCmd.Flags().BoolVarP(&Commit, "commit", "", true, "git commit changes")
	revokeCmd.Flags().BoolVarP(&Push, "push", "", true, "git push changes")
	revokeCmd.Flags().StringVarP(&CertDir, "cert", "c", "cert", "Cert Directory")
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

		if Commit {
			for i := 0; i < len(args); i++ {
				msg := "[infra-admin] Revoke " + args[i] + "certificate"
				files := []string{
					path.Join(CertDir, "crl.pem"),
					path.Join(CertDir, "index.txt"),
					path.Join(CertDir, "index.txt.attr"),
					path.Join(CertDir, "reqs"),
					path.Join(CertDir, "certs_by_serial"),
					path.Join(CertDir, "issued"),
					path.Join(CertDir, "revoked"),
				}
				git.Add(files)
				git.Commit(files, msg)
			}
		}

		if Push {
			git.Push()
		}

	},
}
