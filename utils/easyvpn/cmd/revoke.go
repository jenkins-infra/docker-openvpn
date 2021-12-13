package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/easyrsa"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/git"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/helpers"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/network"
	"github.com/spf13/cobra"
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
		helpers.DecryptPrivateDir()
		errors := easyrsa.RevokeClientCert(args)
		if errors != nil {
			for _, err := range errors {
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}
		for i := range args {
			network.DeleteClientConfig(path.Join(CertDir, "ccd", args[i]))
		}

		fileToDelete := []string{
			path.Join(CertDir, "pki", "index.txt.old"),
			path.Join(CertDir, "pki", "index.txt.attr.old"),
		}

		for i := range fileToDelete {
			if err := os.Remove(fileToDelete[i]); err != nil {
				fmt.Println(err)
			}
		}

		if Commit {
			for i := 0; i < len(args); i++ {
				msg := "[infra-admin] Revoke " + args[i] + " certificate"
				files := []string{
					path.Join(CertDir, "pki", "crl.pem"),
					path.Join(CertDir, "pki", "index.txt"),
					path.Join(CertDir, "pki", "issued", args[i]+".crt"),
					path.Join(CertDir, "pki", "reqs", args[i]+".req"),
					path.Join(CertDir, "pki", "certs_by_serial"),
					path.Join(CertDir, "pki", "index.txt.attr"),
					path.Join(CertDir, "pki", "revoked"),
					path.Join(CertDir, "ccd", args[i]),
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
