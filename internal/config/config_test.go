package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	example := Example()
	if !strings.Contains(example, "endpoints:") {
		t.Error("Example should contain endpoints")
	}
	if !strings.Contains(example, "/etcdhosts") {
		t.Error("Example should contain /etcdhosts")
	}
	t.Logf("Example config:\n%s", example)
}

func TestLoad_Defaults(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(configPath, []byte("endpoints:\n  - http://localhost:2379\n"), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Key != "/etcdhosts" {
		t.Errorf("Default key = %s, want /etcdhosts", cfg.Key)
	}
	if cfg.ReqTimeout != 5*time.Second {
		t.Errorf("Default ReqTimeout = %v, want 5s", cfg.ReqTimeout)
	}
	if cfg.DialTimeout != 5*time.Second {
		t.Errorf("Default DialTimeout = %v, want 5s", cfg.DialTimeout)
	}
}

func TestLoad_Full(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test.yaml")
	content := `endpoints:
  - https://10.0.0.1:2379
  - https://10.0.0.2:2379
key: /custom/key
dial_timeout: 10s
req_timeout: 15s
ca: /path/to/ca.pem
cert: /path/to/cert.pem
cert_key: /path/to/key.pem
username: admin
password: secret
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(cfg.Endpoints) != 2 {
		t.Errorf("Endpoints count = %d, want 2", len(cfg.Endpoints))
	}
	if cfg.Key != "/custom/key" {
		t.Errorf("Key = %s, want /custom/key", cfg.Key)
	}
	if cfg.DialTimeout != 10*time.Second {
		t.Errorf("DialTimeout = %v, want 10s", cfg.DialTimeout)
	}
	if cfg.ReqTimeout != 15*time.Second {
		t.Errorf("ReqTimeout = %v, want 15s", cfg.ReqTimeout)
	}
	if cfg.Username != "admin" {
		t.Errorf("Username = %s, want admin", cfg.Username)
	}
}

func TestLoad_NotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("Load() with nonexistent file should return error")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")
	if err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644); err != nil {
		t.Fatalf("Failed to create config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Load() with invalid YAML should return error")
	}
}

func TestToClientConfig(t *testing.T) {
	cfg := &Config{
		Endpoints:   []string{"http://localhost:2379"},
		Key:         "/test",
		DialTimeout: 5 * time.Second,
		ReqTimeout:  10 * time.Second,
		CA:          "/ca.pem",
		Cert:        "/cert.pem",
		CertKey:     "/key.pem",
		Username:    "user",
		Password:    "pass",
	}

	clientCfg := cfg.ToClientConfig()

	if len(clientCfg.Endpoints) != 1 {
		t.Errorf("Endpoints count = %d, want 1", len(clientCfg.Endpoints))
	}
	if clientCfg.Key != "/test" {
		t.Errorf("Key = %s, want /test", clientCfg.Key)
	}
	if clientCfg.Username != "user" {
		t.Errorf("Username = %s, want user", clientCfg.Username)
	}
}
