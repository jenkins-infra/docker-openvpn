package cmd

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/clientconfig"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/config"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/easyrsa"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/git"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().BoolVarP(&commit, "commit", "", true, "git commit changes")
	signCmd.Flags().BoolVarP(&push, "push", "", true, "git push changes")
	signCmd.Flags().BoolVarP(&cleanup, "cleanup", "", true, "cleanup local sensitive files")
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
		// Load Current Configuration
		globalConfig := config.Config{}
		err := globalConfig.ReadConfigFile(configuration)
		if err != nil {
			log.Fatalf("[ERROR] configuration file %q is invalid: could not read or parse it: %s\n", configuration, err)
		}

		_, ok := globalConfig.Networks[mainNetwork]
		if !ok {
			log.Fatalf("[ERROR] Network %s not found: check config file %s.\n", mainNetwork, configuration)
		}

		// Check if specified users exist in the configuration
		var missingUsers []string
		for _, user := range args {
			_, found := globalConfig.Users[user]
			if !found {
				missingUsers = append(missingUsers, user)
			}
		}
		if len(missingUsers) > 0 {
			log.Fatalf(
				"[ERROR] Missing users from the configuration file %s. You can add them manually or with the help of the command 'easyvpn config %s'.",
				configuration,
				strings.Join(missingUsers, " "),
			)
		}

		// Validate
		// Sign requested certificate(s)
		helpers.DecryptPrivateDir()
		errors := easyrsa.SignClientRequest(args)
		for _, err := range errors {
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}

		for i, username := range args {
			fmt.Printf("Generating Client Configuration for user %s...\n", username)
			err := clientconfig.CreateClientConfig(clientConfigsDir, globalConfig.Users[username], globalConfig.Networks[mainNetwork])
			if err != nil {
				log.Fatalf("[ERROR] could not initialize the client configuration for user %q: %s\n", username, err)
			}
			fmt.Println("Client Configuration added!")

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

		if push {
			git.Push()
		}

		if cleanup {
			helpers.CleanPrivateDir()
		}
	},
}
