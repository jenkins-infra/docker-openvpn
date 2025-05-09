package cmd

import (
	"fmt"
	"log"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/config"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/users"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	configAllRoutes bool
	configRoutes    []string
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVarP(&configAllRoutes, "all-routes", "", false, "Should the user have all routes (e.g. admin user)?")
	configCmd.Flags().StringVarP(&configuration, "config", "c", "config.yaml", "Network Configuration File.")
	configCmd.Flags().StringVarP(&mainNetwork, "net", "", "private", "Network assigned.")
	configCmd.Flags().StringSliceVarP(&configRoutes, "routes", "r", []string{}, "List of custom routes in the specified network.\nMutually exclusive with the flag 'all-routes'.")
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure client network ip",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		globalConfig := &config.Config{}
		err := globalConfig.ReadConfigFile(configuration)
		if err != nil {
			log.Fatalf("[ERROR] configuration file %q is invalid: could not read or parse it: %s\n", configuration, err)
		}

		_, ok := globalConfig.Networks[mainNetwork]
		if !ok {
			log.Fatalf("[ERROR] configuration file %q is invalid: could not find the network %q\n", configuration, mainNetwork)
		}

		if configAllRoutes && len(configRoutes) > 0 {
			log.Fatalln("[ERROR] Both flags 'routes' and 'all-routes' provided, but they are mutually exclusives.")
		}

		for _, routeToCheck := range configRoutes {
			_, ok := globalConfig.Networks[mainNetwork].Routes[routeToCheck]
			if !ok {
				log.Fatalf("[ERROR] The route %q provided by the flag 'routes' does not exist in the configuration file %s for the network %q.", routeToCheck, configuration, mainNetwork)
			}
		}

		userList := args
		newUserConfigs := map[string]users.User{}
		for _, username := range userList {
			newUser := users.User{
				Id:        globalConfig.GetNextUserId(),
				Name:      username,
				AllRoutes: configAllRoutes,
			}
			if len(configRoutes) > 0 {
				newUser.Routes = map[string][]string{
					mainNetwork: configRoutes,
				}
			}
			newUserConfigs[username] = newUser
			yamlContent, err := yaml.Marshal(newUserConfigs)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Add the following YAML snippet in your configuration file %s under the key '$.users':\n\n", configuration)

			fmt.Println("---")
			fmt.Println(string(yamlContent))
			fmt.Println("...")
		}
	},
}
