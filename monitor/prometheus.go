package monitor

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/streasure/util/component"
)

type StatsProvider interface {
	Stats() Stats
}

type Stats struct {
	ConnectionsTotal   uint64
	ConnectionsActive  int64
	ConnectionsCreated uint64
	ConnectionsClosed  uint64

	MessagesReceived  int64
	MessagesForwarded int64
	MessagesPushed    int64
	MessagesPerSecond float64
	MessagesProcessed int64
	MessagesFailed    int64

	MessagesDroppedOverload       int64
	MessagesDroppedFull           int64
	MessagesDroppedNoConn         int64
	MessagesDroppedNoLogic        int64
	MessagesDroppedNoLogicNotConn int64
	MessagesPushDroppedNoConn     int64
	MessagesDroppedBlacklist      int64
	MessagesDroppedRateLimit      int64
	MessagesDroppedWAF            int64
	MessagesDroppedCircuit        int64
	MessagesDroppedIntegrity      int64
	MessagesDroppedFilterChain    int64

	LatencyP50Us        int64
	LatencyP95Us        int64
	LatencyP99Us        int64
	LatencyMaxUs        int64
	ProcessingTimeAvgUs float64

	WAFBlocked            int64
	RateLimitHits         uint64
	RateLimitBlocked      uint64
	CircuitBreakerTripped int64
	DegradationTriggered  int64

	IsLeader int64

	CanaryHit              int64
	TrafficMirrorForwarded int64
	TrafficMirrorDropped   int64

	AlertSent    int64
	AlertDropped int64

	Goroutines      int
	MemoryAlloc     uint64
	MemorySys       uint64
	GCCount         uint32
	CPUUsagePercent float64
	MemUsagePercent float64

	Custom map[string]float64
}

type ExporterConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Addr    string `yaml:"addr" json:"addr"`
	Path    string `yaml:"path" json:"path"`
	Prefix  string `yaml:"prefix" json:"prefix"` // metric prefix, default "app"
}

type GrafanaDatasource struct {
	Name      string `yaml:"name" json:"name"`
	Type      string `yaml:"type" json:"type"`
	Access    string `yaml:"access" json:"access"`
	URL       string `yaml:"url" json:"url"`
	IsDefault bool   `yaml:"isDefault" json:"isDefault"`
	Editable  bool   `yaml:"editable" json:"editable"`
}

type GrafanaDashboard struct {
	Name   string `yaml:"name" json:"name"`
	File   string `yaml:"file" json:"file"`
	Folder string `yaml:"folder" json:"folder"`
}

type GrafanaConfig struct {
	Datasources []GrafanaDatasource `yaml:"datasources" json:"datasources"`
	Dashboards  []GrafanaDashboard  `yaml:"dashboards" json:"dashboards"`
}

type Exporter struct {
	component.BaseComponent
	cfg      ExporterConfig
	prefix   string
	provider StatsProvider
	server   *http.Server
	mu       sync.Mutex
}

func NewExporter(cfg ExporterConfig, provider StatsProvider) *Exporter {
	if cfg.Addr == "" {
		cfg.Addr = ":9100"
	}
	if cfg.Path == "" {
		cfg.Path = "/metrics"
	}
	prefix := cfg.Prefix
	if prefix == "" {
		prefix = "app"
	}
	return &Exporter{
		cfg:      cfg,
		prefix:   prefix,
		provider: provider,
	}
}

func (e *Exporter) Name() string { return "prometheus-exporter" }

func (e *Exporter) Init() error {
	return nil
}

func (e *Exporter) Start() error {
	if !e.cfg.Enabled {
		return nil
	}
	mux := http.NewServeMux()
	mux.HandleFunc(e.cfg.Path, e.serveHTTP)
	e.server = &http.Server{Addr: e.cfg.Addr, Handler: mux}
	go func() {
		if err := e.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("[monitor] prometheus exporter error: %v\n", err)
		}
	}()
	fmt.Printf("[monitor] prometheus exporter started addr=%s path=%s\n", e.cfg.Addr, e.cfg.Path)
	return nil
}

func (e *Exporter) Destroy() {
	if e.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		e.server.Shutdown(ctx)
	}
}

func (e *Exporter) serveHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var s Stats
	if e.provider != nil {
		s = e.provider.Stats()
	} else {
		s = e.systemStats()
	}
	fmt.Fprint(w, RenderPrometheusTextWithPrefix(s, e.prefix))
}

func (e *Exporter) systemStats() Stats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return Stats{
		Goroutines:  runtime.NumGoroutine(),
		MemoryAlloc: m.Alloc,
		MemorySys:   m.Sys,
		GCCount:     m.NumGC,
	}
}

func RenderPrometheusText(s Stats) string {
	return RenderPrometheusTextWithPrefix(s, "app")
}

func RenderPrometheusTextWithPrefix(s Stats, prefix string) string {
	if prefix == "" {
		prefix = "app"
	}
	p := prefix + "_"
	var b strings.Builder

	writeMetric(&b, p+"connections_total", "counter", "Total connections", s.ConnectionsTotal)
	writeMetricI(&b, p+"connections_active", "gauge", "Active connections", s.ConnectionsActive)
	writeMetric(&b, p+"connections_created", "counter", "Created connections", s.ConnectionsCreated)
	writeMetric(&b, p+"connections_closed", "counter", "Closed connections", s.ConnectionsClosed)

	writeMetricI(&b, p+"messages_received_total", "counter", "Received messages", s.MessagesReceived)
	writeMetricI(&b, p+"messages_forwarded_total", "counter", "Forwarded messages", s.MessagesForwarded)
	writeMetricI(&b, p+"messages_pushed_total", "counter", "Pushed messages", s.MessagesPushed)
	writeMetricF(&b, p+"messages_per_second", "gauge", "Messages per second", s.MessagesPerSecond)
	writeMetricI(&b, p+"messages_processed_total", "counter", "Processed messages", s.MessagesProcessed)
	writeMetricI(&b, p+"messages_failed_total", "counter", "Failed messages", s.MessagesFailed)

	writeMetricI(&b, p+"messages_dropped_overload_total", "counter", "Dropped overload", s.MessagesDroppedOverload)
	writeMetricI(&b, p+"messages_dropped_full_total", "counter", "Dropped queue full", s.MessagesDroppedFull)
	writeMetricI(&b, p+"messages_dropped_no_conn_total", "counter", "Dropped no connection", s.MessagesDroppedNoConn)
	writeMetricI(&b, p+"messages_dropped_no_logic_total", "counter", "Dropped no logic server", s.MessagesDroppedNoLogic+s.MessagesDroppedNoLogicNotConn)
	writeMetricI(&b, p+"messages_dropped_push_no_conn_total", "counter", "Push dropped no client connection", s.MessagesPushDroppedNoConn)
	writeMetricI(&b, p+"messages_dropped_blacklist_total", "counter", "Dropped blacklist", s.MessagesDroppedBlacklist)
	writeMetricI(&b, p+"messages_dropped_rate_limit_total", "counter", "Dropped rate limit", s.MessagesDroppedRateLimit)
	writeMetricI(&b, p+"messages_dropped_waf_total", "counter", "Dropped WAF", s.MessagesDroppedWAF)
	writeMetricI(&b, p+"messages_dropped_circuit_total", "counter", "Dropped circuit breaker", s.MessagesDroppedCircuit)
	writeMetricI(&b, p+"messages_dropped_integrity_total", "counter", "Dropped integrity check", s.MessagesDroppedIntegrity)
	writeMetricI(&b, p+"messages_dropped_filter_chain_total", "counter", "Dropped filter chain", s.MessagesDroppedFilterChain)

	writeMetricI(&b, p+"latency_p50_us", "gauge", "P50 latency us", s.LatencyP50Us)
	writeMetricI(&b, p+"latency_p95_us", "gauge", "P95 latency us", s.LatencyP95Us)
	writeMetricI(&b, p+"latency_p99_us", "gauge", "P99 latency us", s.LatencyP99Us)
	writeMetricI(&b, p+"latency_max_us", "gauge", "Max latency us", s.LatencyMaxUs)
	writeMetricF(&b, p+"processing_time_avg_us", "gauge", "Avg processing time us", s.ProcessingTimeAvgUs)

	writeMetricI(&b, p+"waf_blocked_total", "counter", "WAF blocked", s.WAFBlocked)
	writeMetric(&b, p+"rate_limit_hits_total", "counter", "Rate limit hits", s.RateLimitHits)
	writeMetric(&b, p+"rate_limit_blocked_total", "counter", "Rate limit blocked", s.RateLimitBlocked)
	writeMetricI(&b, p+"circuit_breaker_tripped_total", "counter", "Circuit breaker trips", s.CircuitBreakerTripped)
	writeMetricI(&b, p+"degradation_triggered_total", "counter", "Degradation triggers", s.DegradationTriggered)
	writeMetricI(&b, p+"is_leader", "gauge", "Is leader", s.IsLeader)

	writeMetricI(&b, p+"canary_hit_total", "counter", "Canary release hits", s.CanaryHit)
	writeMetricI(&b, p+"traffic_mirror_forwarded_total", "counter", "Mirrored messages forwarded", s.TrafficMirrorForwarded)
	writeMetricI(&b, p+"traffic_mirror_dropped_total", "counter", "Mirrored messages dropped", s.TrafficMirrorDropped)

	writeMetricI(&b, p+"alert_sent_total", "counter", "Alerts sent", s.AlertSent)
	writeMetricI(&b, p+"alert_dropped_total", "counter", "Alerts dropped", s.AlertDropped)

	writeMetricI(&b, p+"goroutines", "gauge", "Goroutines", int64(s.Goroutines))
	writeMetric(&b, p+"memory_alloc_bytes", "gauge", "Memory alloc bytes", s.MemoryAlloc)
	writeMetric(&b, p+"memory_sys_bytes", "gauge", "Memory sys bytes", s.MemorySys)
	writeMetric(&b, p+"gc_count", "counter", "GC count", uint64(s.GCCount))
	writeMetricF(&b, p+"cpu_percent", "gauge", "CPU percent", s.CPUUsagePercent)
	writeMetricF(&b, p+"mem_percent", "gauge", "Mem percent", s.MemUsagePercent)

	for k, v := range s.Custom {
		writeMetricF(&b, p+"custom_"+k, "gauge", "Custom metric: "+k, v)
	}

	return b.String()
}

func writeMetric(b *strings.Builder, name, typ, help string, value uint64) {
	fmt.Fprintf(b, "# HELP %s %s\n# TYPE %s %s\n%s %d\n\n", name, help, name, typ, name, value)
}

func writeMetricI(b *strings.Builder, name, typ, help string, value int64) {
	fmt.Fprintf(b, "# HELP %s %s\n# TYPE %s %s\n%s %d\n\n", name, help, name, typ, name, value)
}

func writeMetricF(b *strings.Builder, name, typ, help string, value float64) {
	fmt.Fprintf(b, "# HELP %s %s\n# TYPE %s %s\n%s %.2f\n\n", name, help, name, typ, name, value)
}
