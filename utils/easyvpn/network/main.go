package network

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Config represents a list of network
type Config struct {
	Networks map[string]Network `yaml:"networks"`
}

// Network represents a network information
type Network struct {
	Name    string
	IPRange string            `yaml:"iprange"`
	NetMask string            `yaml:"netmask"`
	Routes  map[string]string `yaml:"routes"`
}

type clientConfig struct {
	IP      string
	Netmask string
	Routes  []string
}

var clientConfigTemplate = `ifconfig-push {{ .IP }} {{ .Netmask }}
	{{- range .Routes }}
push "route {{ . }}"
	{{- end }}
`

// ReadConfigFile reads a network configuration file
func ReadConfigFile(path string) *Config {
	file, err := ioutil.ReadFile(path)
	config := Config{}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for key, network := range config.Networks {
		network.Name = key
		config.Networks[key] = network
	}

	return &config
}

// Increment IP
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// CheckErr print error essage
func CheckErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

func (n *Network) iprange() ([]string, error) {
	var ips []string
	ip, ipnet, err := net.ParseCIDR(n.IPRange)
	CheckErr(err)

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func isConfigFileExist(cn string) bool {
	if _, err := os.Stat(cn); err != nil {
		return false
	}
	return true
}

func readClientConfigFile(cn string) (ip, mask string) {
	var IP, netmask string
	file, err := os.Open(cn)
	CheckErr(err)

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ifconfig-push") {
			line := strings.Split(line, " ")
			if len(line) == 3 {
				IP = line[1]
				netmask = line[2]
				break
			}
		}
	}
	return IP, netmask
}

func getAllUsedIP(ccd string) ([]string, error) {
	var ipUsed []string
	dir, err := os.Open(ccd)
	CheckErr(err)
	files, err := dir.Readdirnames(-1)
	for i := 0; i < len(files); i++ {
		file := path.Join(ccd, files[i])
		if isConfigFileExist(file) {
			ip, _ := readClientConfigFile(file)
			ipUsed = append(ipUsed, ip)
		}

	}
	return ipUsed, err
}

func isClientConfigured(cn string, clientNetwork string) bool {
	if !isConfigFileExist(cn) {
		return false
	}

	file, err := os.Open(cn)
	defer file.Close()
	CheckErr(err)

	ip, _ := readClientConfigFile(cn)

	_, network, err := net.ParseCIDR(clientNetwork)
	if network.Contains(net.ParseIP(ip)) {
		return true
	}
	return false
}

// DeleteClientConfig remove a client network configuration
func DeleteClientConfig(path string) error {
	err := os.Remove(path)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

// convertRoutesFormat() convert 192.168.0.1/24 CIDR to 192.168.0.1 255.255.255.0
func (n *Network) convertRoutesFormat() []string {
	var result []string
	for _, route := range n.Routes {
		ip, network, err := net.ParseCIDR(route)
		CheckErr(err)
		networkMask := net.IP(network.Mask).String()
		result = append(result, fmt.Sprintf("%v %v", ip.String(), networkMask))
	}
	return result
}

// CreateClientConfig will generate a new client configuration under ccd
func (n *Network) CreateClientConfig(cn string, ccd string) error {
	if isClientConfigured(path.Join(ccd, cn), n.IPRange) {
		fmt.Printf("%v is already in network: %v\n", cn, n.Name)
		return nil
	}

	freeIP, err := n.getFreeIP(ccd)
	CheckErr(err)

	config := clientConfig{
		IP:      freeIP,
		Netmask: net.ParseIP(n.NetMask).String(),
		Routes:  n.convertRoutesFormat(),
	}

	tmpl, err := template.New(cn).Parse(clientConfigTemplate)
	CheckErr(err)

	file, err := os.Create(path.Join(ccd, cn))
	CheckErr(err)
	defer file.Close()

	err = tmpl.Execute(file, config)

	CheckErr(err)
	return err
}

func (n *Network) getFreeIP(ccd string) (string, error) {
	_, network, err := net.ParseCIDR(n.IPRange)

	if err != nil {
		fmt.Println(err)
	}
	networkMask, _ := network.Mask.Size()
	networkIP := network.IP.String()

	networkCIDR := fmt.Sprintf("%v/%v", networkIP, networkMask)

	iprange, err := n.iprange()
	ipUsed, err := getAllUsedIP(ccd)

	for j := 0; j < len(ipUsed); j++ {
		// Restart from 0 as ipUsed is not sorted
		for i := 0; i < len(iprange); i++ {
			if network.Contains(net.ParseIP(ipUsed[j])) && (iprange[i] == ipUsed[j]) {
				fmt.Printf("Found used ip: %v\n", ipUsed[j])
				iprange = append(iprange[:i], iprange[i+1:]...)
				break
			}
		}

	}

	if len(iprange) < 2 {
		msg := fmt.Sprintf("%v doesn't have free ip anymore", networkCIDR)
		err := errors.New(msg)
		fmt.Printf(msg)
		return "", err
	}
	return iprange[1], nil
}
