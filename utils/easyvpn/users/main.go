package users

// User represents a user allowed to access the VPN in the configuration
type User struct {
	Name      string              `yaml:"-"`
	Id        int                 `yaml:"id"`
	AllRoutes bool                `yaml:"all_routes,omitempty"`
	Routes    map[string][]string `yaml:"routes,omitempty"`
}
