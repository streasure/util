package etcd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/streasure/util/component"
	"gopkg.in/yaml.v3"
)

var _ component.Component = (*Component)(nil)

type testConfig struct {
	Name string `json:"name" yaml:"name" toml:"name"`
}

func (testConfig) Validate() error { return nil }

func TestSingleComponentLifecycleAndValidation(t *testing.T) {
	c := New(ComponentConfig{Enabled: false})
	if c.Name() != "etcd" {
		t.Fatalf("unexpected component name: %q", c.Name())
	}
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}
	if err := c.Start(); err != nil {
		t.Fatal(err)
	}
	c.Destroy()
	c.Destroy()
}

func TestRegistrationValidation(t *testing.T) {
	tests := []ComponentConfig{
		{Enabled: true, Registration: RegistrationConfig{Enabled: true, ServiceID: "svc", Address: "127.0.0.1:1"}},
		{Enabled: true, Registration: RegistrationConfig{Enabled: true, ServiceID: "svc", InstanceID: "i", Address: "127.0.0.1:1", LeaseTTL: "500ms"}},
		{Enabled: true, Discovery: DiscoveryConfig{Enabled: true}},
		{Enabled: true, Config: DynamicConfig{Enabled: true, Key: "/cfg", Format: "ini"}},
	}
	for _, cfg := range tests {
		if err := New(cfg).Init(); err == nil {
			t.Fatalf("expected invalid config: %#v", cfg)
		}
	}
}

func TestDecodeConfigFormats(t *testing.T) {
	for _, test := range []struct{ format, data string }{
		{"json", `{"name":"json"}`},
		{"yaml", "name: yaml\n"},
		{"toml", "name = \"toml\"\n"},
	} {
		var value testConfig
		if err := DecodeConfig([]byte(test.data), test.format, &value); err != nil {
			t.Fatalf("%s: %v", test.format, err)
		}
		if value.Name != test.format {
			t.Fatalf("%s decoded as %q", test.format, value.Name)
		}
	}
	if err := validateConfigData([]byte("{"), "json"); err == nil {
		t.Fatal("invalid JSON should fail")
	}
}

func TestLoadBalanceAndServiceSet(t *testing.T) {
	c := New(ComponentConfig{Enabled: true, Discovery: DiscoveryConfig{Enabled: true, ServiceID: "svc", LoadBalance: LoadBalanceRoundRobin}})
	c.services.Store([]string{"10.0.0.2:80", "10.0.0.1:80"})
	if got, ok := c.Pick(); !ok || got == "" {
		t.Fatal("round robin did not select an endpoint")
	}
	c.applyServiceSet(map[string]string{"/services/svc/b": "10.0.0.2:80", "/services/svc/a": "10.0.0.1:80"})
	services := c.Services()
	if len(services) != 2 || services[0] != "10.0.0.1:80" {
		t.Fatalf("services are not deterministic: %#v", services)
	}
}

func TestEtcdConfigFile(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "config", "etcd.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	var config ComponentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatal(err)
	}
	if err := New(config).Init(); err != nil {
		t.Fatal(err)
	}
	if endpoints := config.Etcd.endpointList(); len(endpoints) != 1 || endpoints[0] != "http://127.0.0.1:2379" {
		t.Fatalf("unexpected endpoints: %#v", endpoints)
	}
}

func TestTLSConfigValidation(t *testing.T) {
	if err := (Config{Endpoint: "http://127.0.0.1:2379", CACertFile: "ca.pem"}).validate(); err == nil {
		t.Fatal("TLS on an HTTP endpoint should be rejected")
	}
	if err := (Config{Endpoint: "https://127.0.0.1:2379", CertFile: "client.pem"}).validate(); err == nil {
		t.Fatal("incomplete client certificate settings should be rejected")
	}
}
