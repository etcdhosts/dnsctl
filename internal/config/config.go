package config

import (
	"os"
	"time"

	client "github.com/etcdhosts/client-go/v2"
	"gopkg.in/yaml.v3"
)

// Config holds the dnsctl configuration.
type Config struct {
	Endpoints   []string      `yaml:"endpoints"`
	Key         string        `yaml:"key,omitempty"`
	DialTimeout time.Duration `yaml:"dial_timeout,omitempty"`
	ReqTimeout  time.Duration `yaml:"req_timeout,omitempty"`
	CA          string        `yaml:"ca,omitempty"`
	Cert        string        `yaml:"cert,omitempty"`
	CertKey     string        `yaml:"cert_key,omitempty"`
	Username    string        `yaml:"username,omitempty"`
	Password    string        `yaml:"password,omitempty"`
}

// ToClientConfig converts Config to client.Config.
func (c *Config) ToClientConfig() client.Config {
	return client.Config{
		Endpoints:   c.Endpoints,
		Key:         c.Key,
		DialTimeout: c.DialTimeout,
		ReqTimeout:  c.ReqTimeout,
		CA:          c.CA,
		Cert:        c.Cert,
		CertKey:     c.CertKey,
		Username:    c.Username,
		Password:    c.Password,
	}
}

// Load loads config from file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Apply defaults
	if cfg.Key == "" {
		cfg.Key = "/etcdhosts"
	}
	if cfg.ReqTimeout == 0 {
		cfg.ReqTimeout = 5 * time.Second
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5 * time.Second
	}

	return &cfg, nil
}

// Example returns an example config YAML string.
func Example() string {
	cfg := Config{
		Endpoints: []string{
			"https://172.16.1.21:2379",
			"https://172.16.1.22:2379",
			"https://172.16.1.23:2379",
		},
		Key:     "/etcdhosts",
		CA:      "/etc/etcd/ssl/etcd-ca.pem",
		Cert:    "/etc/etcd/ssl/etcd.pem",
		CertKey: "/etc/etcd/ssl/etcd-key.pem",
	}
	out, _ := yaml.Marshal(cfg)
	return string(out)
}
