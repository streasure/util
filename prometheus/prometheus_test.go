package prometheus

import (
	"strings"
	"testing"
)

type mockStatsProvider struct{}

func (m *mockStatsProvider) Stats() Stats {
	return Stats{
		ConnectionsTotal:  100,
		ConnectionsActive: 10,
		MessagesReceived:  5000,
		MessagesProcessed: 4900,
		MessagesFailed:    100,
		Goroutines:        20,
		MemoryAlloc:       1024 * 1024,
		MemorySys:         2 * 1024 * 1024,
		GCCount:           5,
		Custom: map[string]float64{
			"test_metric": 42.0,
		},
	}
}

func TestRenderPrometheusText(t *testing.T) {
	s := Stats{
		ConnectionsTotal: 100,
		MessagesReceived: 5000,
		Goroutines:       20,
	}
	text := RenderPrometheusText(s)

	if !strings.Contains(text, "app_connections_total 100") {
		t.Error("missing connections_total metric")
	}
	if !strings.Contains(text, "app_messages_received_total 5000") {
		t.Error("missing messages_received_total metric")
	}
	if !strings.Contains(text, "app_goroutines 20") {
		t.Error("missing goroutines metric")
	}
	if !strings.Contains(text, "# HELP") {
		t.Error("missing HELP annotation")
	}
	if !strings.Contains(text, "# TYPE") {
		t.Error("missing TYPE annotation")
	}
}

func TestRenderCustomMetrics(t *testing.T) {
	s := Stats{
		Custom: map[string]float64{
			"my_counter": 99.5,
		},
	}
	text := RenderPrometheusText(s)
	if !strings.Contains(text, "app_custom_my_counter 99.50") {
		t.Error("missing custom metric")
	}
}

func TestExporterDisabled(t *testing.T) {
	cfg := ExporterConfig{Enabled: false}
	exp := NewExporter(cfg, nil)
	if err := exp.Start(); err != nil {
		t.Errorf("Start on disabled exporter should not error: %v", err)
	}
	exp.Destroy()
}

func TestExporterDefaults(t *testing.T) {
	cfg := ExporterConfig{Enabled: true, Addr: "", Path: ""}
	exp := NewExporter(cfg, nil)
	if exp.cfg.Addr != ":9100" {
		t.Errorf("default addr = %q, want :9100", exp.cfg.Addr)
	}
	if exp.cfg.Path != "/metrics" {
		t.Errorf("default path = %q, want /metrics", exp.cfg.Path)
	}
}

func TestMockStatsProvider(t *testing.T) {
	p := &mockStatsProvider{}
	s := p.Stats()
	if s.ConnectionsTotal != 100 {
		t.Errorf("ConnectionsTotal = %d, want 100", s.ConnectionsTotal)
	}
	if s.Custom["test_metric"] != 42.0 {
		t.Errorf("Custom[test_metric] = %f, want 42.0", s.Custom["test_metric"])
	}
}
