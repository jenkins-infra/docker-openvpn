package network

// Network represents a network information
type Network struct {
	Name    string            `yaml:"-"`
	IPRange string            `yaml:"iprange"`
	Routes  map[string]string `yaml:"routes"`
}
