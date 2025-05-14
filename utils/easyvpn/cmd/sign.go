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
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().BoolVarP(&commit, "commit", "", true, "git commit changes")
	signCmd.Flags().BoolVarP(&push, "push", "", true, "git push changes")
	signCmd.Flags().StringVarP(&certDir, "certsDir", "", "cert", "Cert Directory")
	signCmd.Flags().StringVarP(&clientConfigsDir, "ccd", "", "cert/ccd", "Client Config Directory")
	signCmd.Flags().StringVarP(&configuration, "config", "", "config.yaml", "Network Configuration File")
	signCmd.Flags().StringVarP(&mainNetwork, "net", "n", "private", "Network to assign the cn")
}

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a request certificate",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Sign requested certificate(s)
		helpers.DecryptPrivateDir()
		errors := easyrsa.SignClientRequest(args)
		for _, err := range errors {
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}

		// Generate client config
		globalConfig := network.ReadConfigFile(configuration)

		network, ok := globalConfig.Networks[mainNetwork]
		if !ok {
			fmt.Printf("Network %s not found: check config file %s.\n", mainNetwork, configuration)
			os.Exit(1)
		}

		for i := range args {
			err := network.CreateClientConfig(args[i], path.Join(clientConfigsDir, mainNetwork))
			if err != nil {
				panic(err)
			}

			// Commit changes
			if commit {
				msg := "[infra-admin] Sign certificate request for " + args[i]
				files := []string{
					path.Join(certDir, "pki/issued", args[i]+".crt"),
					path.Join(certDir, "pki", "index.txt"),
					path.Join(certDir, "pki", "index.txt.attr"),
					path.Join(certDir, "pki", "certs_by_serial"),
					path.Join(certDir, "pki", "serial"),
					path.Join(certDir, "ccd", mainNetwork, args[i]),
				}
				git.Add(files)
				git.Commit(files, msg)
			}
		}

		// Push changes
		if push {
			git.Push()
		}

		helpers.CleanPrivateDir()
	},
}
