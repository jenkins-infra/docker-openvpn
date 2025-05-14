package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/easyrsa"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/git"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(revokeCmd)
	revokeCmd.Flags().BoolVarP(&commit, "commit", "", true, "git commit changes")
	revokeCmd.Flags().BoolVarP(&push, "push", "", true, "git push changes")
	revokeCmd.Flags().StringVarP(&certDir, "cert", "c", "cert", "Cert Directory")
	revokeCmd.Flags().StringVarP(&mainNetwork, "network", "n", "private", "mainNetwork")
}

var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke a client certificate",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		helpers.DecryptPrivateDir()
		errors := easyrsa.RevokeClientCert(args)
		for _, err := range errors {
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
		for i := range args {
			err := os.Remove(path.Join(certDir, "ccd", mainNetwork, args[i]))
			if err != nil {
				fmt.Println(err)
				fmt.Println("Continuing despite error...")
			}
		}

		fileToDelete := []string{
			path.Join(certDir, "pki", "index.txt.old"),
			path.Join(certDir, "pki", "index.txt.attr.old"),
		}

		for i := range fileToDelete {
			if err := os.Remove(fileToDelete[i]); err != nil {
				fmt.Println(err)
				fmt.Println("Continuing despite error...")
			}
		}

		if commit {
			for i := range args {
				msg := "[infra-admin] Revoke " + args[i] + " certificate"
				files := []string{
					path.Join(certDir, "pki", "crl.pem"),
					path.Join(certDir, "pki", "index.txt"),
					path.Join(certDir, "pki", "issued", args[i]+".crt"),
					path.Join(certDir, "pki", "reqs", args[i]+".req"),
					path.Join(certDir, "pki", "certs_by_serial"),
					path.Join(certDir, "pki", "index.txt.attr"),
					path.Join(certDir, "pki", "revoked"),
					path.Join(certDir, "ccd", mainNetwork, args[i]),
				}
				git.Add(files)
				git.Commit(files, msg)
			}
		}

		if push {
			git.Push()
		}

		helpers.CleanPrivateDir()
	},
}
