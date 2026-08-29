package nacos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/streasure/util/component"
)

type DiscoveryConfig struct {
	Enabled      bool         `yaml:"enabled" json:"enabled"`
	Nacos        Config       `yaml:"nacos" json:"nacos"`
	Service      NamingConfig `yaml:"service" json:"service"`
	ScanInterval string       `yaml:"scanInterval" json:"scanInterval"`
	Zone         string       `yaml:"zone" json:"zone"`
}

type ServiceChangeCallback func(event ServiceEvent)

type Discovery struct {
	component.BaseComponent
	cfg       DiscoveryConfig
	client    *client
	services  map[string]*ServiceInfo
	mu        sync.RWMutex
	callbacks []ServiceChangeCallback
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

func NewDiscovery(cfg DiscoveryConfig) *Discovery {
	return &Discovery{
		cfg:       cfg,
		client:    newClient(cfg.Nacos),
		services:  make(map[string]*ServiceInfo),
		callbacks: make([]ServiceChangeCallback, 0),
		stopCh:    make(chan struct{}),
	}
}

func (d *Discovery) Name() string { return "nacos-discovery" }

func (d *Discovery) Init() error { return nil }

func (d *Discovery) Start() error {
	if !d.cfg.Enabled {
		return nil
	}
	if d.cfg.Nacos.Endpoint == "" {
		return fmt.Errorf("nacos endpoint empty, discovery disabled")
	}
	if err := d.pullServices(); err != nil {
		fmt.Printf("[nacos] discovery initial pull failed: %v\n", err)
	}
	d.wg.Add(1)
	go d.scanLoop()
	fmt.Printf("[nacos] discovery started service=%s interval=%s\n", d.cfg.Service.ServiceName, d.cfg.ScanInterval)
	return nil
}

func (d *Discovery) Destroy() {
	if !d.cfg.Enabled {
		return
	}
	close(d.stopCh)
	d.wg.Wait()
}

func (d *Discovery) OnServiceChange(cb ServiceChangeCallback) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.callbacks = append(d.callbacks, cb)
}

func (d *Discovery) GetServices() []*ServiceInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]*ServiceInfo, 0, len(d.services))
	for _, svc := range d.services {
		result = append(result, svc)
	}
	return result
}

func (d *Discovery) GetService(serviceID string) *ServiceInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.services[serviceID]
}

func (d *Discovery) GetServiceByAddress(address string) *ServiceInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, svc := range d.services {
		if svc.Address == address {
			return svc
		}
	}
	return nil
}

func (d *Discovery) scanLoop() {
	defer d.wg.Done()
	interval := pollInterval(d.cfg.ScanInterval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-d.stopCh:
			return
		case <-ticker.C:
			if err := d.pullServices(); err != nil {
				fmt.Printf("[nacos] discovery pull error: %v\n", err)
			}
		}
	}
}

type nacosInstance struct {
	InstanceID string            `json:"instanceId"`
	IP         string            `json:"ip"`
	Port       int               `json:"port"`
	Weight     float64           `json:"weight"`
	Metadata   map[string]string `json:"metadata"`
	Healthy    bool              `json:"healthy"`
	Enabled    bool              `json:"enabled"`
}

func (d *Discovery) pullServices() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, _ := d.client.ensureToken(ctx)

	q := url.Values{}
	q.Set("serviceName", d.cfg.Service.ServiceName)
	q.Set("groupName", d.client.cfg.Group)
	q.Set("namespaceId", d.client.cfg.Namespace)

	var reqURL string
	if strings.ToLower(d.client.cfg.APIVersion) == "v1" {
		reqURL = fmt.Sprintf("%s/nacos/v1/ns/instance/list?%s", d.client.cfg.Endpoint, q.Encode())
	} else {
		reqURL = fmt.Sprintf("%s/nacos/v3/client/ns/instance/list?%s", d.client.namingEndpoint(), q.Encode())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := d.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("nacos list status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	instances, err := d.parseInstances(body)
	if err != nil {
		return err
	}
	d.reconcile(instances)
	return nil
}

func (d *Discovery) parseInstances(body []byte) ([]*ServiceInfo, error) {
	var wrapper struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	var instances []nacosInstance

	if err := json.Unmarshal(body, &wrapper); err == nil && len(wrapper.Data) > 0 {
		if arrErr := json.Unmarshal(wrapper.Data, &instances); arrErr != nil || instances == nil {
			var paged struct {
				Instances []nacosInstance `json:"instances"`
				PageItems []nacosInstance `json:"pageItems"`
			}
			if err := json.Unmarshal(wrapper.Data, &paged); err != nil {
				return nil, fmt.Errorf("unmarshal data: %w", err)
			}
			instances = append(instances, paged.Instances...)
			instances = append(instances, paged.PageItems...)
		}
	} else {
		var data struct {
			Instances []nacosInstance `json:"instances"`
		}
		if err := json.Unmarshal(body, &data); err != nil {
			return nil, fmt.Errorf("unmarshal v1: %w", err)
		}
		instances = data.Instances
	}

	result := make([]*ServiceInfo, 0, len(instances))
	for _, inst := range instances {
		if !inst.Healthy || !inst.Enabled {
			continue
		}
		id := inst.InstanceID
		if id == "" {
			id = fmt.Sprintf("%s:%d", inst.IP, inst.Port)
		}
		svc := &ServiceInfo{
			ServiceID:   id,
			ServiceName: d.cfg.Service.ServiceName,
			Address:     fmt.Sprintf("%s:%d", inst.IP, inst.Port),
			Weight:      int(inst.Weight),
			Metadata:    inst.Metadata,
		}
		if d.cfg.Zone != "" {
			svcZone := svc.Metadata["zone"]
			if svcZone == "" {
				svcZone = "default"
			}
			if svcZone != d.cfg.Zone {
				continue
			}
		}
		result = append(result, svc)
	}
	return result, nil
}

func (d *Discovery) reconcile(newServices []*ServiceInfo) {
	activeMap := make(map[string]*ServiceInfo, len(newServices))
	for _, svc := range newServices {
		activeMap[svc.ServiceID] = svc
	}

	d.mu.Lock()
	oldServices := make(map[string]*ServiceInfo, len(d.services))
	for k, v := range d.services {
		oldServices[k] = v
	}
	d.services = make(map[string]*ServiceInfo, len(newServices))
	for id, svc := range activeMap {
		d.services[id] = svc
	}
	d.mu.Unlock()

	for id, svc := range activeMap {
		if _, existed := oldServices[id]; !existed {
			d.notifyCallbacks(ServiceEvent{
				Type:      EventRegister,
				Service:   *svc,
				Timestamp: time.Now().UnixMilli(),
			})
		}
	}
	for id, svc := range oldServices {
		if _, exists := activeMap[id]; !exists {
			d.notifyCallbacks(ServiceEvent{
				Type:      EventDeregister,
				Service:   *svc,
				Timestamp: time.Now().UnixMilli(),
			})
		}
	}
}

func (d *Discovery) notifyCallbacks(event ServiceEvent) {
	d.mu.RLock()
	cbs := make([]ServiceChangeCallback, len(d.callbacks))
	copy(cbs, d.callbacks)
	d.mu.RUnlock()
	for _, cb := range cbs {
		func() {
			defer func() { recover() }()
			cb(event)
		}()
	}
}
