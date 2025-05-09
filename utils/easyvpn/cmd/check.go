package cmd

import (
	"os"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/checks"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&certDir, "cert", "c", "cert", "Cert Directory")
	checkCmd.Flags().StringVarP(&mainNetwork, "net", "", "private", "Network assigned")
	checkCmd.Flags().StringVarP(&configuration, "config", "f", "config.yaml", "Network Configuration File")
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks various repository config",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		rc := 0

		globalConfig := &config.Config{}
		err := globalConfig.ReadConfigFile(configuration)
		if err != nil {
			rc = 1
		}

		if result, _ := checks.IsAllCertsSigned(certDir); !result {
			rc = 1
		}
		if result, _ := checks.IsAllClientConfigured(certDir, mainNetwork, *globalConfig); !result {
			rc = 1
		}

		os.Exit(rc)
	},
}
