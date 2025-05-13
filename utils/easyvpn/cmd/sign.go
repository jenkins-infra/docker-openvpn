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

var certDir string
var commit bool
var config string
var push bool

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().BoolVarP(&commit, "commit", "", true, "git commit changes")
	signCmd.Flags().BoolVarP(&push, "push", "", true, "git push changes")
	signCmd.Flags().StringVarP(&certDir, "certsDir", "", "cert", "Cert Directory")
	signCmd.Flags().StringVarP(&ccd, "ccd", "", "cert/ccd", "Client Config Directory")
	signCmd.Flags().StringVarP(&config, "config", "", "config.yaml", "Network Configuration File")
	signCmd.Flags().StringVarP(&net, "net", "n", "private", "Network to assign the cn")

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
		globalConfig := network.ReadConfigFile(config)

		network, ok := globalConfig.Networks[net]
		if !ok {
			fmt.Printf("Network %s not found: check config file %s.\n", net, config)
			os.Exit(1)
		}

		for i := range args {
			err := network.CreateClientConfig(args[i], path.Join(ccd, net))
			if err != nil {
				panic(err)
			}

			// Commit changes
			if commit {
				msg := "[infra-admin] Sign certificate request for " + args[i]
				files := []string{
					path.Join(CertDir, "pki/issued", args[i]+".crt"),
					path.Join(CertDir, "pki", "index.txt"),
					path.Join(CertDir, "pki", "index.txt.attr"),
					path.Join(CertDir, "pki", "certs_by_serial"),
					path.Join(CertDir, "pki", "serial"),
					path.Join(CertDir, "ccd", net, args[i]),
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
