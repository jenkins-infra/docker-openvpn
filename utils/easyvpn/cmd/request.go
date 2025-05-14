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
	requestCmd.Flags().BoolVarP(&commit, "commit", "", true, "git commit changes")
	requestCmd.Flags().BoolVarP(&push, "push", "", true, "git push changes")
	requestCmd.Flags().StringVarP(&certDir, "cert", "c", "cert", "Cert Directory")
}

var requestCmd = &cobra.Command{
	Use:   "request",
	Short: "Generate a client private key and a request certificate",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		errors := easyrsa.RequestClientCert(args)
		for _, err := range errors {
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
		if commit {
			for i := range args {
				msg := "[infra-admin] Submit certificate request for " + args[i]
				files := []string{
					path.Join(certDir, "pki/reqs", args[i]+".req"),
				}
				git.Add(files)
				git.Commit(files, msg)
			}
		}

		if push {
			git.Push()
		}
	},
}
