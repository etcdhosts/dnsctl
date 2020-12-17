package main

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/etcdhosts/client-go"
	"gopkg.in/yaml.v2"
)

type Config struct {
	configPath    string
	client.Config `yaml:",inline"`
}

func (cfg *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawConfig Config
	raw := rawConfig{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	if raw.HostsKey == "" {
		raw.HostsKey = "/etcdhosts"
	}
	if raw.ReqTimeout == 0 {
		raw.ReqTimeout = 2 * time.Second
	}
	if raw.DialTimeout == 0 {
		raw.DialTimeout = 5 * time.Second
	}

	*cfg = Config(raw)
	return nil
}

func (cfg *Config) SetConfigPath(configPath string) {
	cfg.configPath = configPath
}

func (cfg *Config) Write() error {
	if cfg.configPath == "" {
		return errors.New("config path not set")
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cfg.configPath, out, 0644)
}

func (cfg *Config) WriteTo(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	cfg.configPath = filePath
	return cfg.Write()
}

func (cfg *Config) Load() error {
	if cfg.configPath == "" {
		return errors.New("config path not set")
	}
	buf, err := ioutil.ReadFile(cfg.configPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, cfg)
}

func (cfg *Config) LoadFrom(filePath string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}
	cfg.configPath = filePath
	return cfg.Load()
}

func ExampleConfig() string {
	out, _ := yaml.Marshal(Config{
		Config: client.Config{
			CA:   "/etc/etcd/ssl/etcd-root-ca.pem",
			Cert: "/etc/etcd/ssl/etcd.pem",
			Key:  "/etc/etcd/ssl/etcd-key.pem",
			Endpoints: []string{
				"https://172.16.1.21:2379",
				"https://172.16.1.22:2379",
				"https://172.16.1.23:2379",
			},
		},
	})
	return string(out)
}
