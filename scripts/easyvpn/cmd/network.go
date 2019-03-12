package cmd

import (
	"../network"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

// Delete define if a client conffiguration must be removed
var Delete bool

// Member is the CN assigned to a network
var Member string

// CCD represent the client configuration directory path
var CCD string

// Config contains Network Configuration
var Config string

func init() {
	rootCmd.AddCommand(networkCmd)
	networkCmd.Flags().StringVarP(&Member, "member", "m", "", "Network member")
	networkCmd.Flags().StringVarP(&CCD, "ccd", "", "cert/ccd", "Client Config Directory")
	networkCmd.Flags().StringVarP(&Config, "config", "c", "config.yaml", "Network Configuration File")
	networkCmd.Flags().BoolVarP(&Delete, "delete", "d", false, "Delete Network Configuration File")
	networkCmd.MarkFlagRequired("member")
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Create client configuraton for a specific network",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if Delete == true {
			network.DeleteClientConfig(path.Join(CCD, Member))
			os.Exit(0)
		}

		file, err := ioutil.ReadFile(Config)
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

		networkName := args[0]

		for i := 0; i < len(config.Networks); i++ {
			if config.Networks[i].Name == networkName {
				network.CreateClientConfig(Member, config.Networks[i].CIDR, CCD)
				os.Exit(0)
			}
		}
		fmt.Printf("Network configuration for %s not found in %s\n", networkName, Config)
		os.Exit(1)
	},
}
