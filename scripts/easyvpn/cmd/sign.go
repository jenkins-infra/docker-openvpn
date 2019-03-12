package cmd

import (
	"../easyrsa"
	"../git"
	"../helpers"
	"../network"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
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
	signCmd.Flags().StringVarP(&net, "net", "", "default", "Network to assign the cn")

}

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a request certificate",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Sign requested certificate(s)
		helpers.DecryptPrivateDir()
		errors := easyrsa.SignClientRequest(args)
		if errors != nil {
			for _, err := range errors {
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		}

		// Generate client config
		file, err := ioutil.ReadFile(config)
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

		networkNotFound := false
		for i := 0; i < len(config.Networks); i++ {
			if config.Networks[i].Name == net {
				networkNotFound = true
				for j := 0; j < len(args); j++ {
					network.CreateClientConfig(args[j], config.Networks[i].CIDR, ccd)
				}
			}
		}

		if !networkNotFound {
			fmt.Printf("Network %v not founded", net)
			os.Exit(1)
		}

		// Commit changes
		if commit {
			for i := 0; i < len(args); i++ {
				msg := "[infra-admin] Sign certificate request for " + args[i]
				files := []string{
					path.Join(CertDir, "pki/issued", args[i]+".crt"),
					path.Join(CertDir, "pki", "index.txt"),
					path.Join(CertDir, "pki", "serial"),
					path.Join(CertDir, "ccd", args[i]),
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
