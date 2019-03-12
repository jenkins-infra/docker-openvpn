package cmd

import (
	"../git"
	"../network"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
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
				if push {
					git.Push()
				}
			}
			os.Exit(0)
		} else {

			file, err := ioutil.ReadFile(configuration)
			config := network.Config{}

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = yaml.Unmarshal(file, &config)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			networkFounded := false
			for i := 0; i < len(config.Networks); i++ {
				if config.Networks[i].Name == net {
					networkFounded = true
					for j := 0; j < len(args); j++ {
						network.CreateClientConfig(args[j], config.Networks[i].CIDR, ccd)
						if commit {
							msg := fmt.Sprintf("[infra-admin] Update %v network configuration", args[j])
							files := []string{
								path.Join(ccd, args[j]),
							}
							git.Add(files)
							git.Commit(files, msg)
						}
						if push {
							git.Push()
						}
					}
				}
			}
			if !networkFounded {
				fmt.Printf("Network %v not founded in %v", net, configuration)
				os.Exit(1)
			}
		}
	},
}
