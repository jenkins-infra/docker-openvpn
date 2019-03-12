package network

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"text/template"
)

// Config contains...
type Config struct {
	Networks []struct {
		Name  string   `yaml:"name"`
		CIDR  string   `yaml:"cidr"`
		Mask  string   `yaml:"mask"`
		Route []string `yaml:"routes"`
	} `yaml:"networks"`
}

type clientConfig struct {
	IP      string
	Netmask string
	Routes  []string
	Members []string
}

// DeleteClientConfig
// GetAllNetworkMembers(name string)
//

var clientConfigTemplate = `ifconfig-push {{ .IP }} {{ .Netmask }}
	{{- range .Routes }}
push "route {{ . }}"
	{{- end }}
`

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

func iprange(cidr string) ([]string, error) {
	var ips []string
	ip, ipnet, err := net.ParseCIDR(cidr)
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

// DeleteClientConfig remote a specify network client configuration
func DeleteClientConfig(path string) error {
	err := os.Remove(path)

	if err != nil {
		fmt.Println(err)
	}
	return err
}

// CreateClientConfig will generate a new client configuration under ccd
func CreateClientConfig(cn string, clientNetwork string, ccd string) error {
	if isClientConfigured(path.Join(ccd, cn), clientNetwork) {
		fmt.Printf("Woop woop nothing to do for %v, he is already in the right network \n", cn)
		return nil
	}

	_, network, err := net.ParseCIDR(clientNetwork)

	CheckErr(err)
	freeIP, err := getFreeIP(clientNetwork, ccd)
	CheckErr(err)

	config := clientConfig{
		IP:      freeIP,
		Netmask: net.IP(network.Mask).String(),
		Routes: []string{
			"10.8.0.1 255.255.255.0",
			"10.9.0.1 255.255.255.0",
		},
	}
	tmpl, err := template.New(cn).Parse(clientConfigTemplate)
	CheckErr(err)

	file, err := os.Create(ccd + "/" + cn)
	CheckErr(err)
	defer file.Close()

	err = tmpl.Execute(file, config)

	CheckErr(err)
	return err
}

func getFreeIP(clientNetwork string, ccd string) (string, error) {
	_, network, err := net.ParseCIDR(clientNetwork)

	if err != nil {
		fmt.Println(err)
	}
	networkMask, _ := network.Mask.Size()
	networkIP := network.IP.String()

	networkCIDR := fmt.Sprintf("%v/%v", networkIP, networkMask)

	fmt.Printf("%v/%v\n", networkIP, networkMask)

	iprange, err := iprange(networkCIDR)
	ipUsed, err := getAllUsedIP(ccd)

	for i := 0; i < len(iprange); i++ {
		for j := 0; j < len(ipUsed); j++ {
			if network.Contains(net.ParseIP(ipUsed[j])) && (iprange[i] == ipUsed[j]) {
				fmt.Printf("Found used ip: %v\n", ipUsed[j])
				iprange = append(iprange[:i], iprange[i+1:]...)
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
