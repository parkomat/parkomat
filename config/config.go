package config

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"strings"
)

type Domain struct {
	Name string `json:"name" toml:"name"`

	// TODO: Not implemented yet
	LetsEncrypt bool `json:"lets_encrypt" toml:"lets_encrypt"`

	HasSSL bool `json:"has_ssl" toml:"has_ssl"`

	SSLCertificate    string `json:"ssl_certificate" toml:"ssl_certificate"`
	SSLCertificateKey string `json:"ssl_certificate_key" toml:"ssl_certificate_key"`

	Zone *Zone `json:"zone" toml:"zone"`
}

type Server struct {
	Name string `json:"name" toml:"name"`
	IP   string `json:"ip" toml:"ip"`
}

type Web struct {
	IP        string `json:"ip" toml:"ip"`
	Port      int    `json:"port" toml:"port"`
	SSLPort   int    `json:"ssl_port" toml:"ssl_port"`
	Path      string `json:"path" toml:"path"`
	AccessLog string `json:"access_log" toml:"access_log"`
	ErrorLog  string `json:"error_log" toml:"error_log"`
}

type DNS struct {
	IP      string   `json:"ip" toml:"ip"`
	Port    int      `json:"port" toml:"port"`
	Servers []Server `json:"servers"`
}

type WebDav struct {
	Enabled  bool   `json:"enabled" toml:"enabled"`
	Username string `json:"username" toml:"username"`
	Password string `json:"password" toml:"password"`
	Mount    string `json:"mount" toml:"mount"`
}

type Zone struct {
	A   string `json:"A" toml:"A"`
	MX  string `json:"MX" toml:"MX"`
	TXT string `json:"TXT" toml:"TXT"`
}

type Config struct {
	CatchAll bool `json:"catch_all" toml:"catch_all"`

	Domains []*Domain `json:"domains" toml:"domains"`

	Zone   Zone   `json:"zone" toml:"zone"`
	Web    Web    `json:"web" toml:"web"`
	DNS    DNS    `json:"dns" toml:"dns"`
	WebDav WebDav `json:"webdav" toml:"webdav"`
}

func (config *Config) Refresh() {
	// TODO: not implemented
	// This function should refresh configuration
	// Without the need of restarting the service
}

// NewConfigFromFile reads configuration from specified filename
func NewConfigFromFile(filename string) (config *Config, err error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	if _, err := toml.Decode(string(file), &config); err != nil {
		return nil, err
	}

	if len(config.DNS.Servers) == 0 {
		return nil, fmt.Errorf("No DNS servers specified")
	}

	config.normalizeData()
	return
}

// normalizeData makes sure all domains are in the same format
func (config *Config) normalizeData() {
	for i, _ := range config.Domains {
		d := config.Domains[i]
		d.Name = strings.Trim(strings.ToLower(d.Name), " ")
	}
}

// NewConfigFromJSONFile reads configuration from JSON formatted file
func NewConfigFromJSONFile(filename string) (config *Config, err error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &config)

	if err != nil {
		return nil, err
	}

	config.normalizeData()
	return
}

// HasDomain checks if domain exists in the loaded configuration
func (c *Config) HasDomain(domain string) bool {
	if c.CatchAll == true {
		return true
	}

	for _, d := range c.Domains {
		if strings.Contains(domain, d.Name) {
			return true
		}
	}
	return false
}

// GetDomain gets domain details from the loaded configuration. If CatchAll is enabled it will create new
// Domain record if it doesn't exist
func (c *Config) GetDomain(domain string) *Domain {
	lcDomain := strings.ToLower(domain)
	for i := range c.Domains {
		if strings.Contains(lcDomain, strings.ToLower(c.Domains[i].Name)) {
			return c.Domains[i]
		}
	}

	if c.CatchAll == true {
		// No domain found, add it
		d := &Domain{
			Name: domain,
		}
		c.Domains = append(c.Domains, d)
		return d
	}
	return nil
}
