package nacos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/streasure/util/component"
)

type RegistryConfig struct {
	Enabled           bool          `yaml:"enabled" json:"enabled"`
	Nacos             Config        `yaml:"nacos" json:"nacos"`
	Service           NamingConfig  `yaml:"service" json:"service"`
	HeartbeatInterval time.Duration `yaml:"heartbeatInterval" json:"heartbeatInterval"`
	HeartbeatTTL      time.Duration `yaml:"heartbeatTTL" json:"heartbeatTTL"`
}

type Registry struct {
	component.BaseComponent
	cfg      RegistryConfig
	client   *client
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

func NewRegistry(cfg RegistryConfig) *Registry {
	if cfg.HeartbeatInterval <= 0 {
		cfg.HeartbeatInterval = 3 * time.Second
	}
	if cfg.HeartbeatTTL <= 0 {
		cfg.HeartbeatTTL = 15 * time.Second
	}
	return &Registry{
		cfg:    cfg,
		client: newClient(cfg.Nacos),
		stopCh: make(chan struct{}),
	}
}

func (r *Registry) Name() string { return "nacos-registry" }

func (r *Registry) Init() error {
	return nil
}

func (r *Registry) Start() error {
	if !r.cfg.Enabled {
		return nil
	}
	if r.cfg.Nacos.Endpoint == "" {
		return nil
	}
	if err := r.registerInstance(); err != nil {
		return fmt.Errorf("nacos register: %w", err)
	}
	r.wg.Add(1)
	go r.heartbeatLoop()
	fmt.Printf("[nacos] registry started service=%s addr=%s\n", r.cfg.Service.ServiceName, r.cfg.Service.Addr)
	return nil
}

func (r *Registry) Destroy() {
	if !r.cfg.Enabled {
		return
	}
	close(r.stopCh)
	r.wg.Wait()
	_ = r.deregisterInstance()
}

func (r *Registry) heartbeatLoop() {
	defer r.wg.Done()
	ticker := time.NewTicker(r.cfg.HeartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-r.stopCh:
			return
		case <-ticker.C:
			if err := r.registerInstance(); err != nil {
				fmt.Printf("[nacos] registry heartbeat error: %v\n", err)
			}
		}
	}
}

func (r *Registry) registerInstance() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, _ := r.client.ensureToken(ctx)
	ip, port := parseAddrPort(r.cfg.Service.Addr)
	metadataJSON, _ := json.Marshal(r.cfg.Service.Metadata)

	form := url.Values{}
	form.Set("serviceName", r.cfg.Service.ServiceName)
	form.Set("ip", ip)
	form.Set("port", port)
	form.Set("weight", fmt.Sprintf("%d", r.cfg.Service.Weight))
	form.Set("metadata", string(metadataJSON))
	form.Set("clusterName", "DEFAULT")
	form.Set("groupName", r.client.cfg.Group)
	form.Set("namespaceId", r.client.cfg.Namespace)
	form.Set("ephemeral", "true")

	var reqURL string
	if strings.ToLower(r.client.cfg.APIVersion) == "v1" {
		reqURL = fmt.Sprintf("%s/nacos/v1/ns/instance", r.client.cfg.Endpoint)
	} else {
		reqURL = fmt.Sprintf("%s/nacos/v3/client/ns/instance", r.client.namingEndpoint())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nacos register status %d", resp.StatusCode)
	}
	return nil
}

func (r *Registry) deregisterInstance() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, _ := r.client.ensureToken(ctx)
	ip, port := parseAddrPort(r.cfg.Service.Addr)

	q := url.Values{}
	q.Set("serviceName", r.cfg.Service.ServiceName)
	q.Set("ip", ip)
	q.Set("port", port)
	q.Set("groupName", r.client.cfg.Group)
	q.Set("namespaceId", r.client.cfg.Namespace)
	q.Set("clusterName", "DEFAULT")

	var reqURL string
	if strings.ToLower(r.client.cfg.APIVersion) == "v1" {
		reqURL = fmt.Sprintf("%s/nacos/v1/ns/instance?%s", r.client.cfg.Endpoint, q.Encode())
	} else {
		reqURL = fmt.Sprintf("%s/nacos/v3/client/ns/instance?%s", r.client.namingEndpoint(), q.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
