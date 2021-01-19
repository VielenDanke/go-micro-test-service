package config

// Server ...
type Server struct {
	Name    string `json:"name"`
	Version string `json:"-"`
	Addr    string `json:"addr" default:":9090"`
}

// Config ...
type Config struct {
	Server *Server
}

// NewConfig ...
func NewConfig(name, version string) *Config {
	return &Config{
		Server: &Server{
			Name:    name,
			Version: version,
		},
	}
}
