package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/git"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/network"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&clientConfigsDir, "ccd", "", "cert/ccd", "Client Config Directory")
	configCmd.Flags().StringVarP(&configuration, "config", "c", "config.yaml", "Network Configuration File")
	configCmd.Flags().StringVarP(&mainNetwork, "net", "", "private", "Network assigned")
	configCmd.Flags().BoolVarP(&commit, "commit", "", true, "Commit changes")
	configCmd.Flags().BoolVarP(&push, "push", "", true, "Push changes")
	configCmd.Flags().BoolVarP(&delete, "delete", "d", false, "Delete Network Configuration File")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure client network ip",
	Args:  cobra.MinimumNArgs(1),
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
			globalConfig := network.ReadConfigFile(clientConfigsDir)

			network, ok := globalConfig.Networks[mainNetwork]
			if !ok {
				fmt.Printf("Network %s not found: check config file %s.\n", mainNetwork, configuration)
				os.Exit(1)
			}

			for j := range args {
				user := args[j]
				fmt.Printf("Generating CCD configuration for user %s\n", user)
				err := network.CreateClientConfig(user, path.Join(clientConfigsDir, mainNetwork))
				if err != nil {
					panic(err)
				}

				if commit {
					msg := fmt.Sprintf("[infra-admin] Update %v in '%v' network configuration", user, mainNetwork)
					files := []string{
						path.Join(clientConfigsDir, mainNetwork, user),
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
