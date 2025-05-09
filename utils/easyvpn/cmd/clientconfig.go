package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/clientconfig"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/config"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/git"
	"github.com/spf13/cobra"
)

var clientConfigAllUsers bool

func init() {
	rootCmd.AddCommand(clientConfigCmd)
	clientConfigCmd.Flags().StringVarP(&clientConfigsDir, "ccd", "", "cert/ccd", "Client Config Directory")
	clientConfigCmd.Flags().StringVarP(&configuration, "config", "c", "config.yaml", "Network Configuration File")
	clientConfigCmd.Flags().StringVarP(&mainNetwork, "net", "", "private", "Network assigned")
	clientConfigCmd.Flags().BoolVarP(&commit, "commit", "", true, "Commit changes")
	clientConfigCmd.Flags().BoolVarP(&push, "push", "", true, "Push changes")
	clientConfigCmd.Flags().BoolVarP(&delete, "delete", "d", false, "Delete Network Configuration File")
	clientConfigCmd.Flags().BoolVarP(&clientConfigAllUsers, "all", "a", false, "Generate client configuration for all users with a signed certificate")
}

var clientConfigCmd = &cobra.Command{
	Use:   "clientconfig",
	Short: "Manage Client Configurations",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if delete {
			for j := range args {
				err := os.Remove(path.Join(clientConfigsDir, mainNetwork, args[j]))
				if err != nil {
					fmt.Println(err)
					fmt.Println("Continuing despite error...")
				}
				if commit {
					msg := fmt.Sprintf("[infra-admin] Delete %v in '%v' network configuration", args[j], mainNetwork)
					files := []string{
						path.Join(clientConfigsDir, args[j]),
					}
					git.Add(files)
					git.Commit(files, msg)
				}
			}
		} else {
			globalConfig := &config.Config{}
			err := globalConfig.ReadConfigFile(configuration)
			if err != nil {
				log.Fatalf("[ERROR] configuration file %q is invalid: could not read or parse it: %s\n", configuration, err)
			}

			_, ok := globalConfig.Networks[mainNetwork]
			if !ok {
				log.Fatalf("[ERROR] configuration file %q is invalid: could not find the network %q\n", configuration, mainNetwork)
			}

			userList := args
			if clientConfigAllUsers {
				userList = config.GetUsersWithCertificate(certDir, *globalConfig)
			} else {
				// Validate user provided arguments
				if len(userList) == 0 {
					log.Fatalf("[ERROR] The 'clientconfig' command expects at least 1 argument when the flag --all is not specified.")
				}
				for _, user := range userList {
					_, ok := globalConfig.Users[user]
					if !ok {
						log.Fatalf("[ERROR] The provided user %q is not configured in the file %q. Pleas fix.", user, configuration)
					}
				}
			}

			for _, username := range userList {
				fmt.Printf("Generating Client Configuration for user %s...\n", username)
				err := clientconfig.CreateClientConfig(clientConfigsDir, globalConfig.Users[username], globalConfig.Networks[mainNetwork])
				if err != nil {
					log.Fatalf("[ERROR] could not initialize the client configuration for user %q: %s\n", username, err)
				}
				fmt.Println("Client Configuration added!")

				if commit {
					msg := fmt.Sprintf("[infra-admin] Update %v in '%v' network configuration", username, mainNetwork)
					files := []string{
						path.Join(clientConfigsDir, mainNetwork, username),
					}
					git.Add(files)
					git.Commit(files, msg)
				}
			}
		}

		if push {
			git.Push()
		}
	},
}
