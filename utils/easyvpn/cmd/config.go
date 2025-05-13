package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/git"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/network"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&ccd, "ccd", "", "cert/ccd", "Client Config Directory")
	configCmd.Flags().StringVarP(&configuration, "config", "c", "config.yaml", "Network Configuration File")
	configCmd.Flags().StringVarP(&net, "net", "", "private", "Network assigned")
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
				err := network.DeleteClientConfig(path.Join(ccd, net, args[j]))
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}

				if commit {
					msg := fmt.Sprintf("[infra-admin] Delete %v in '%v' network configuration", args[j], net)
					files := []string{
						path.Join(ccd, args[j]),
					}
					git.Add(files)
					git.Commit(files, msg)
				}
			}
		} else {
			globalConfig := network.ReadConfigFile(config)

			network, ok := globalConfig.Networks[net]
			if !ok {
				fmt.Printf("Network %s not found: check config file %s.\n", net, config)
				os.Exit(1)
			}

			for j := range args {
				user := args[j]
				fmt.Printf("Generating CCD configuration for user %s\n", user)
				err := network.CreateClientConfig(user, path.Join(ccd, net))
				if err != nil {
					panic(err)
				}

				if commit {
					fmt.Println("KO")
					msg := fmt.Sprintf("[infra-admin] Update %v in '%v' network configuration", user, net)
					files := []string{
						path.Join(ccd, net, user),
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
