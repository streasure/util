package nacos

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/streasure/util/component"
)

type ConfigCenterConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Nacos   Config `yaml:"nacos" json:"nacos"`
}

type ConfigChangeCallback func(data []byte)

type ConfigCenter struct {
	component.BaseComponent
	cfg       ConfigCenterConfig
	client    *client
	callbacks []ConfigChangeCallback
	stopCh    chan struct{}
	wg        sync.WaitGroup
	lastHash  string
	mu        sync.RWMutex
}

func NewConfigCenter(cfg ConfigCenterConfig) *ConfigCenter {
	return &ConfigCenter{
		cfg:       cfg,
		client:    newClient(cfg.Nacos),
		callbacks: make([]ConfigChangeCallback, 0),
		stopCh:    make(chan struct{}),
	}
}

func (cc *ConfigCenter) Name() string { return "nacos-config-center" }

func (cc *ConfigCenter) Init() error { return nil }

func (cc *ConfigCenter) Start() error {
	if !cc.cfg.Enabled {
		return nil
	}
	if cc.cfg.Nacos.Endpoint == "" {
		return fmt.Errorf("nacos endpoint empty, config center disabled")
	}
	data, err := cc.Pull()
	if err != nil {
		fmt.Printf("[nacos] config center initial pull failed: %v\n", err)
	} else if len(data) > 0 {
		cc.notifyCallbacks(data)
	}
	cc.wg.Add(1)
	go cc.watchLoop()
	fmt.Printf("[nacos] config center started dataID=%s group=%s\n", cc.cfg.Nacos.DataID, cc.cfg.Nacos.Group)
	return nil
}

func (cc *ConfigCenter) Destroy() {
	if !cc.cfg.Enabled {
		return
	}
	close(cc.stopCh)
	cc.wg.Wait()
}

func (cc *ConfigCenter) OnConfigChange(cb ConfigChangeCallback) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc.callbacks = append(cc.callbacks, cb)
}

func (cc *ConfigCenter) Pull() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()

	token, _ := cc.client.ensureToken(ctx)
	pullURL := cc.buildPullURL()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pullURL, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := cc.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("config center pull status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	content := cc.unwrap(body)

	cc.mu.Lock()
	defer cc.mu.Unlock()
	h := hashContent(content)
	if cc.lastHash == h {
		return content, nil
	}
	cc.lastHash = h
	return content, nil
}

func (cc *ConfigCenter) watchLoop() {
	defer cc.wg.Done()
	interval := pollInterval(cc.cfg.Nacos.PollInterval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-cc.stopCh:
			return
		case <-ticker.C:
			data, err := cc.Pull()
			if err != nil {
				continue
			}
			if len(data) > 0 {
				cc.notifyCallbacks(data)
			}
		}
	}
}

func (cc *ConfigCenter) notifyCallbacks(data []byte) {
	cc.mu.RLock()
	cbs := make([]ConfigChangeCallback, len(cc.callbacks))
	copy(cbs, cc.callbacks)
	cc.mu.RUnlock()
	for _, cb := range cbs {
		func() {
			defer func() { recover() }()
			cb(data)
		}()
	}
}

func (cc *ConfigCenter) buildPullURL() string {
	nacos := cc.cfg.Nacos
	group := nacos.Group
	if group == "" {
		group = "DEFAULT_GROUP"
	}
	if strings.ToLower(nacos.APIVersion) == "v1" {
		return fmt.Sprintf("%s/nacos/v1/cs/configs?dataId=%s&group=%s&tenant=%s",
			nacos.Endpoint, nacos.DataID, group, nacos.Namespace)
	}
	return fmt.Sprintf("%s/v3/console/cs/config?dataId=%s&groupName=%s&namespaceId=%s",
		nacos.Endpoint, nacos.DataID, group, nacos.Namespace)
}

func (cc *ConfigCenter) unwrap(body []byte) []byte {
	if strings.ToLower(cc.cfg.Nacos.APIVersion) == "v1" {
		return body
	}
	var resp struct {
		Code int    `json:"code"`
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return body
	}
	if resp.Code != 0 {
		return nil
	}
	if resp.Data.Content == "" {
		return nil
	}
	return []byte(resp.Data.Content)
}

func hashContent(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	const prime = 1099511628211
	h := uint64(14695981039346656037)
	for _, c := range b {
		h ^= uint64(c)
		h *= prime
	}
	return fmt.Sprintf("%x", h)
}

func base64Decode(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}
