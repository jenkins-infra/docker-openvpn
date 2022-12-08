package cmd

import (
	"os"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/checks"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().StringVarP(&CertDir, "cert", "c", "cert", "Cert Directory")
	checkCmd.Flags().StringVarP(&Network, "network", "net", "default", "Network")
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks various repository config",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		rc := 0

		if result, _ := checks.IsAllCertsSigned(CertDir); !result {
			rc = 1
		}
		if result, _ := checks.IsAllClientConfigured(CertDir, Network); !result {
			rc = 1
		}

		os.Exit(rc)
	},
}
