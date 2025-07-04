package clientconfig

import (
	"fmt"
	"html/template"
	"net"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/network"
	"github.com/jenkins-infra/docker-openvpn/utils/easyvpn/users"
)

type ClientConfig struct {
	Name        string
	IP          string
	Netmask     string
	NetworkName string
	Routes      []Route
}

type Route struct {
	Name    string
	IP      string
	Netmask string
}

var clientConfigTemplate = `# Client Configuration for user "{{ .Name }}" in the VPN network "{{ .NetworkName }}"
ifconfig-push {{ .IP }} {{ .Netmask }}
	{{- range .Routes }}
# {{ .Name }}{{ if ne .Netmask "255.255.255.255" }} vnet{{ end }}
push "route {{ .IP }} {{ .Netmask }}"
	{{- end }}
`

func initClientConfig(userConfig users.User, networkConfig network.Network) (ClientConfig, error) {
	cc := ClientConfig{
		Name:        userConfig.Name,
		NetworkName: networkConfig.Name,
	}

	// Determine the client IPv4 in the VPN network and the associated IPv4 netmask
	clientIp, clientIpNetwork, err := net.ParseCIDR(networkConfig.IPRange)
	if err != nil {
		return cc, err
	}

	// Note: it's a crud but easy to maintain technique given the low amount of users we have
	splitClientIp := strings.Split(clientIp.String(), ".")
	splitClientIp[3] = fmt.Sprintf("%d", userConfig.Id)
	cc.IP = strings.Join(splitClientIp, ".")
	cc.Netmask = net.IP(clientIpNetwork.Mask).String()

	// Detects the list of routes to add
	routes := make(map[string]string)
	if userConfig.AllRoutes {
		routes = networkConfig.Routes
	} else {
		for _, userRoute := range userConfig.Routes[networkConfig.Name] {
			routes[userRoute] = networkConfig.Routes[userRoute]
		}
	}

	// Sort keys to ensure deterministic output
	keys := make([]string, 0, len(routes))
	for k := range routes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Add routes to the CC
	for _, routeConfigName := range keys {
		routes := strings.Split(routes[routeConfigName], " ")

		for i, routeCidr := range routes {
			routeIP, routeNetwork, err := net.ParseCIDR(routeCidr)
			if err != nil {
				return cc, err
			}
			routeName := routeConfigName
			if len(routes) > 1 {
				routeName = fmt.Sprintf("%s_%d", routeName, (i + 1))
			}
			userRoute := Route{
				Name:    routeName,
				IP:      routeIP.String(),
				Netmask: net.IP(routeNetwork.Mask).String(),
			}
			cc.Routes = append(cc.Routes, userRoute)
		}
	}

	return cc, nil
}

// CreateClientConfig generates a new client config file in the provided ccd folder
func CreateClientConfig(ccd string, userConfig users.User, networkConfig network.Network) error {
	cc, err := initClientConfig(userConfig, networkConfig)
	if err != nil {
		return fmt.Errorf("[ERROR] could not initialize the client configuration for user %q: %s", userConfig.Name, err)
	}

	err = cc.createClientConfig(ccd)
	if err != nil {
		return fmt.Errorf("[ERROR] could not generate client configuration file for user %q: %s", userConfig.Name, err)
	}

	return nil
}

func (cc *ClientConfig) createClientConfig(ccd string) error {
	tmpl, err := template.New(cc.Name).Parse(clientConfigTemplate)
	if err != nil {
		return err
	}

	ccFile := path.Join(ccd, cc.NetworkName, cc.Name)
	file, err := os.Create(ccFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, cc)
	if err != nil {
		return err
	}

	err = file.Close()

	return err
}
