package cmd

import (
	"fmt"
	"path"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/easyrsa"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(requestCmd)
	requestCmd.Flags().BoolVarP(&Commit, "commit", "", true, "git commit changes")
	requestCmd.Flags().BoolVarP(&Push, "push", "", true, "git push changes")
	requestCmd.Flags().StringVarP(&CertDir, "cert", "c", "cert", "Cert Directory")
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
		if Commit {
			for i := 0; i < len(args); i++ {
				msg := "[infra-admin] Submit certificate request for " + args[i]
				files := []string{
					path.Join(CertDir, "pki/reqs", args[i]+".req"),
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
