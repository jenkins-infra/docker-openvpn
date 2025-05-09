package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/helpers"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/network"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/users"
	"gopkg.in/yaml.v2"
)

// Config represents the configuration object read from YAML configuration file
type Config struct {
	Networks map[string]network.Network `yaml:"networks"`
	Users    map[string]users.User      `yaml:"users"`
}

// ReadConfigFile reads configuration from the specified file
func (c *Config) ReadConfigFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, c)
	if err != nil {
		return err
	}

	// Set up network names from their map key
	for key, network := range c.Networks {
		network.Name = key
		c.Networks[key] = network
	}

	// Set up user names from their map key
	for key, user := range c.Users {
		user.Name = key
		c.Users[key] = user

		// Validate specified use routes
		for networkId, routes := range user.Routes {
			for _, routeId := range routes {
				_, ok := c.Networks[networkId].Routes[routeId]
				if !ok {
					log.Fatalf("user %q has an invalid route %q specified for the VPN network %q", user.Name, routeId, networkId)
				}
			}
		}
	}

	return nil
}

// GetUsers returns the sorted list of User names from the configuration
func (c *Config) GetUsers() (userList []string) {
	for userId := range c.Users {
		userList = append(userList, userId)
	}
	slices.Sort(userList)

	return userList
}

func GetUsersWithCertificate(certDir string, globalConfig Config) (userList []string) {
	issuedCertificates := helpers.GetUsernameFile(path.Join(certDir, "pki/issued"), ".crt")
	configUserList := globalConfig.GetUsers()

	// TODO: can we manage server-side somewhere else than the client CRL and certs?
	// And remove the server-side certificate from this list (not a user!)
	for _, cert := range issuedCertificates {
		if cert != "vpn.jenkins.io.crt" {
			userList = append(userList, strings.TrimSuffix(cert, ".crt"))
		}
	}
	slices.Sort(userList)

	if !slices.Equal(configUserList, userList) {
		fmt.Println("[WARNING] there is a mismatch between the list of configured users and the issued certificates.")
	}

	return userList
}

func (c *Config) GetNextUserId() int {
	var currentUserIds []int
	for _, user := range c.Users {
		currentUserIds = append(currentUserIds, user.Id)
	}

	// Required for the following iteration (from smaller to bigger)
	slices.Sort(currentUserIds)

	// The ID 0 and 1 cannot be used (broadcast IP and gateway IP respectively) so we start at 2
	startId := 2

	// If there is a hole in the (sorted) list of IDs, then return
	for i, userId := range currentUserIds {
		if userId < startId {
			log.Fatalf("[ERROR] found a user with an ID < %d which is not allowed. Please fix your configuration.\n", startId)
		}
		nextId := i + startId

		if userId != nextId {
			return nextId
		}
	}

	return len(currentUserIds) + startId
}
