package main

type ServerConfig struct {
	Name    string `json:"name"`
	Version string `json:"-"`
	Addr    string `json:"addr" default:":9090"`
}

type Config struct {
	Server *ServerConfig
}

func NewConfig(name, version string) *Config {
	return &Config{
		Server: &ServerConfig{
			Name:    name,
			Version: version,
		},
	}
}
