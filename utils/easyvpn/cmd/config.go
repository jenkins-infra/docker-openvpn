package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/git"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/network"
	"github.com/spf13/cobra"
)

var delete bool
var ccd string
var configuration string
var net string

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringVarP(&ccd, "ccd", "", "cert/ccd", "Client Config Directory")
	configCmd.Flags().StringVarP(&configuration, "config", "c", "config.yaml", "Network Configuration File")
	configCmd.Flags().StringVarP(&net, "net", "", "default", "Network assigned")
	configCmd.Flags().BoolVarP(&commit, "commit", "", true, "Commit changes")
	configCmd.Flags().BoolVarP(&commit, "push", "", true, "Push changes")
	configCmd.Flags().BoolVarP(&delete, "delete", "d", false, "Delete Network Configuration File")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure client network ip",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if delete == true {
			for j := 0; j < len(args); j++ {
				network.DeleteClientConfig(path.Join(ccd, args[j]))
				if commit {
					msg := fmt.Sprintf("[infra-admin] Delete %v network configuration", args[j])
					files := []string{
						path.Join(ccd, args[j]),
					}
					git.Add(files)
					git.Commit(files, msg)
				}
			}
		} else {

			config := network.ReadConfigFile(config)
			network, err := config.GetNetworkByName(net)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			for j := 0; j < len(args); j++ {
				network.CreateClientConfig(args[j], ccd)

				if commit {
					msg := fmt.Sprintf("[infra-admin] Update %v network configuration", args[j])
					files := []string{
						path.Join(ccd, args[j]),
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
