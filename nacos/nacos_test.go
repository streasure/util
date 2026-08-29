package nacos

import (
	"testing"
	"time"
)

func TestParseAddrPort(t *testing.T) {
	tests := []struct {
		addr string
		ip   string
		port string
	}{
		{"127.0.0.1:8080", "127.0.0.1", "8080"},
		{"10.0.0.1:9090", "10.0.0.1", "9090"},
		{"localhost", "localhost", "0"},
		{"[::1]:8080", "[::1]", "8080"},
	}
	for _, tt := range tests {
		ip, port := parseAddrPort(tt.addr)
		if ip != tt.ip || port != tt.port {
			t.Errorf("parseAddrPort(%q) = (%q, %q), want (%q, %q)", tt.addr, ip, port, tt.ip, tt.port)
		}
	}
}

func TestPollInterval(t *testing.T) {
	if d := pollInterval(""); d != 5*time.Second {
		t.Errorf("empty = %v, want 5s", d)
	}
	if d := pollInterval("10s"); d != 10*time.Second {
		t.Errorf("10s = %v", d)
	}
	if d := pollInterval("100ms"); d != time.Second {
		t.Errorf("100ms = %v, want 1s (clamped)", d)
	}
	if d := pollInterval("0s"); d != time.Second {
		t.Errorf("0s = %v, want 1s", d)
	}
	if d := pollInterval("invalid"); d != 5*time.Second {
		t.Errorf("invalid = %v, want 5s", d)
	}
}

func TestClientDefaults(t *testing.T) {
	c := newClient(Config{})
	if c.cfg.Group != "DEFAULT_GROUP" {
		t.Errorf("default group = %q", c.cfg.Group)
	}
	if c.cfg.APIVersion != "v3" {
		t.Errorf("default apiVersion = %q", c.cfg.APIVersion)
	}
}

func TestRegistryDisabled(t *testing.T) {
	cfg := RegistryConfig{Enabled: false}
	r := NewRegistry(cfg)
	if err := r.Start(); err != nil {
		t.Errorf("Start on disabled should not error: %v", err)
	}
	r.Destroy()
}

func TestDiscoveryDisabled(t *testing.T) {
	cfg := DiscoveryConfig{Enabled: false}
	d := NewDiscovery(cfg)
	if err := d.Start(); err != nil {
		t.Errorf("Start on disabled should not error: %v", err)
	}
	d.Destroy()
}

func TestConfigCenterDisabled(t *testing.T) {
	cfg := ConfigCenterConfig{Enabled: false}
	cc := NewConfigCenter(cfg)
	if err := cc.Start(); err != nil {
		t.Errorf("Start on disabled should not error: %v", err)
	}
	cc.Destroy()
}

func TestDiscoveryGetServicesEmpty(t *testing.T) {
	d := NewDiscovery(DiscoveryConfig{})
	svcs := d.GetServices()
	if len(svcs) != 0 {
		t.Error("empty discovery should return no services")
	}
}

func TestDiscoveryGetServiceNotFound(t *testing.T) {
	d := NewDiscovery(DiscoveryConfig{})
	if svc := d.GetService("nonexistent"); svc != nil {
		t.Error("should return nil for nonexistent service")
	}
}

func TestDiscoveryGetServiceByAddressNotFound(t *testing.T) {
	d := NewDiscovery(DiscoveryConfig{})
	if svc := d.GetServiceByAddress("1.2.3.4:5"); svc != nil {
		t.Error("should return nil for nonexistent address")
	}
}

func TestServiceInfo(t *testing.T) {
	svc := ServiceInfo{
		ServiceID:   "id-1",
		ServiceName: "test-service",
		Address:     "127.0.0.1:8080",
		Weight:      1,
		Metadata:    map[string]string{"zone": "default"},
	}
	if svc.ServiceID != "id-1" {
		t.Error("ServiceID mismatch")
	}
	if svc.Address != "127.0.0.1:8080" {
		t.Error("Address mismatch")
	}
}
