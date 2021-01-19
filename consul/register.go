package consul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var registeredServicesID []string
var client *http.Client
var c *config

func init() {
	client = &http.Client{}
	c = &config{}
}

type config struct {
	address string
}

// Service ...
type Service struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
	Port int      `json:"port"`
}

// SetupConfig ...
func SetupConfig(address string) {
	c.address = address
}

// NewService ...
func NewService(id string, name string, port int) *Service {
	return &Service{
		ID:   id,
		Name: name,
		Port: port,
	}
}

// Register ...
func Register(service *Service) error {
	if c.address == "" {
		return fmt.Errorf("Address of consul should not be empty: %s", c.address)
	}
	data, err := json.Marshal(service)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/agent/service/register", c.address), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	registeredServicesID = append(registeredServicesID, service.ID)
	return nil
}

// Unregister ...
func Unregister() error {
	for _, v := range registeredServicesID {
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/agent/service/deregister/%s", c.address, v), nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			return err
		}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	return nil
}
