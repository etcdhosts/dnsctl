//go:build integration

package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	client "github.com/etcdhosts/client-go/v2"
	"github.com/etcdhosts/dnsctl/v2/internal/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// startEtcd starts an etcd container and returns the endpoint.
func startEtcd(t *testing.T) (string, func()) {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "quay.io/coreos/etcd:v3.5.17",
		ExposedPorts: []string{"2379/tcp"},
		Cmd: []string{
			"/usr/local/bin/etcd",
			"--name", "etcd0",
			"--listen-client-urls", "http://0.0.0.0:2379",
			"--advertise-client-urls", "http://0.0.0.0:2379",
		},
		WaitingFor: wait.ForLog("ready to serve client requests").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start etcd container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "2379")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, port.Port())

	cleanup := func() {
		_ = container.Terminate(ctx)
	}

	return endpoint, cleanup
}

// dumpKeys prints all keys under a prefix from etcd.
func dumpKeys(t *testing.T, endpoint, prefix string) {
	t.Helper()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create etcd client: %v", err)
	}
	defer func() { _ = cli.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		t.Fatalf("Failed to get keys: %v", err)
	}

	t.Logf("\n=== etcd key layout (prefix: %s) ===", prefix)
	for _, kv := range resp.Kvs {
		value := string(kv.Value)
		if len(value) > 80 {
			value = value[:77] + "..."
		}
		value = strings.ReplaceAll(value, "\n", "\\n")
		t.Logf("  %s = %s", string(kv.Key), value)
	}
	t.Logf("=== total: %d keys ===\n", len(resp.Kvs))
}

// createTestConfig creates a temporary config file for testing.
func createTestConfig(t *testing.T, endpoint, key string) string {
	t.Helper()
	content := fmt.Sprintf(`endpoints:
  - %s
key: %s
dial_timeout: 5s
req_timeout: 5s
`, endpoint, key)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "dnsctl.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	return configPath
}

func TestIntegration_AddListDelPurge(t *testing.T) {
	endpoint, cleanup := startEtcd(t)
	defer cleanup()

	configPath := createTestConfig(t, endpoint, "/etcdhosts")

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	cli, err := client.NewClient(cfg.ToClientConfig())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer func() { _ = cli.Close() }()

	t.Run("Add", func(t *testing.T) {
		hosts, _ := cli.Read()
		_ = hosts.Add(client.Record{Hostname: "example.com", IP: net.ParseIP("192.168.1.1"), Weight: 1})
		_ = hosts.Add(client.Record{Hostname: "example.com", IP: net.ParseIP("192.168.1.2"), Weight: 2})
		_ = hosts.Add(client.Record{Hostname: "api.example.com", IP: net.ParseIP("10.0.0.1"), TTL: 60})
		if err := cli.Write(hosts); err != nil {
			t.Fatalf("Write error = %v", err)
		}

		hosts2, _ := cli.Read()
		if hosts2.Len() != 3 {
			t.Errorf("Len = %d, want 3", hosts2.Len())
		}

		dumpKeys(t, endpoint, "/etcdhosts")
	})

	t.Run("List", func(t *testing.T) {
		hosts, err := cli.Read()
		if err != nil {
			t.Fatalf("Read error = %v", err)
		}

		entries := hosts.LookupV4("example.com")
		if len(entries) != 2 {
			t.Errorf("LookupV4(example.com) = %d, want 2", len(entries))
		}

		entries = hosts.LookupV4("api.example.com")
		if len(entries) != 1 {
			t.Errorf("LookupV4(api.example.com) = %d, want 1", len(entries))
		}
		if entries[0].TTL != 60 {
			t.Errorf("TTL = %d, want 60", entries[0].TTL)
		}
	})

	t.Run("Del", func(t *testing.T) {
		hosts, _ := cli.Read()
		if err := hosts.Del("example.com", net.ParseIP("192.168.1.1")); err != nil {
			t.Fatalf("Del error = %v", err)
		}
		if err := cli.Write(hosts); err != nil {
			t.Fatalf("Write error = %v", err)
		}

		hosts2, _ := cli.Read()
		entries := hosts2.LookupV4("example.com")
		if len(entries) != 1 {
			t.Errorf("After Del, LookupV4(example.com) = %d, want 1", len(entries))
		}
		if entries[0].IP.String() != "192.168.1.2" {
			t.Errorf("Remaining IP = %s, want 192.168.1.2", entries[0].IP.String())
		}

		dumpKeys(t, endpoint, "/etcdhosts")
	})

	t.Run("Purge", func(t *testing.T) {
		hosts, _ := cli.Read()
		hosts.Purge("example.com")
		if err := cli.Write(hosts); err != nil {
			t.Fatalf("Write error = %v", err)
		}

		hosts2, _ := cli.Read()
		entries := hosts2.LookupV4("example.com")
		if len(entries) != 0 {
			t.Errorf("After Purge, LookupV4(example.com) = %d, want 0", len(entries))
		}

		entries = hosts2.LookupV4("api.example.com")
		if len(entries) != 1 {
			t.Errorf("After Purge, LookupV4(api.example.com) = %d, want 1", len(entries))
		}

		dumpKeys(t, endpoint, "/etcdhosts")
	})
}

func TestIntegration_KeyLayout(t *testing.T) {
	endpoint, cleanup := startEtcd(t)
	defer cleanup()

	configPath := createTestConfig(t, endpoint, "/dns/records")

	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	cli, err := client.NewClient(cfg.ToClientConfig())
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	defer func() { _ = cli.Close() }()

	hosts, _ := cli.Read()
	_ = hosts.Add(client.Record{Hostname: "web.local", IP: net.ParseIP("10.0.0.1")})
	_ = hosts.Add(client.Record{Hostname: "db.local", IP: net.ParseIP("10.0.0.2"), TTL: 300})
	_ = hosts.Add(client.Record{Hostname: "cache.local", IP: net.ParseIP("10.0.0.3"), Weight: 5})
	_ = cli.Write(hosts)

	t.Log("Custom key prefix /dns/records:")
	dumpKeys(t, endpoint, "/dns")

	mode, err := cli.Mode()
	if err != nil {
		t.Fatalf("Mode() error = %v", err)
	}
	if mode != client.ModeSingle {
		t.Errorf("Mode = %s, want single", mode)
	}

	if cli.Key() != "/dns/records" {
		t.Errorf("Key = %s, want /dns/records", cli.Key())
	}
}
