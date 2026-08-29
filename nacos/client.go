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
)

type Config struct {
	Enabled        bool   `yaml:"enabled" json:"enabled"`
	Type           string `yaml:"type" json:"type"`
	Endpoint       string `yaml:"endpoint" json:"endpoint"`
	NamingEndpoint string `yaml:"namingEndpoint" json:"namingEndpoint"`
	Namespace      string `yaml:"namespace" json:"namespace"`
	DataID         string `yaml:"dataID" json:"dataID"`
	Group          string `yaml:"group" json:"group"`
	Token          string `yaml:"token" json:"token"`
	Username       string `yaml:"username" json:"username"`
	Password       string `yaml:"password" json:"password"`
	PollInterval   string `yaml:"pollInterval" json:"pollInterval"`
	APIVersion     string `yaml:"apiVersion" json:"apiVersion"`
}

type NamingConfig struct {
	ServiceName string `yaml:"serviceName" json:"serviceName"`
	Addr        string `yaml:"addr" json:"addr"`
	Weight      int    `yaml:"weight" json:"weight"`
	Zone        string `yaml:"zone" json:"zone"`
	Metadata    map[string]string `yaml:"metadata" json:"metadata"`
}

type ServiceInfo struct {
	ServiceID   string            `json:"service_id"`
	ServiceName string            `json:"service_name"`
	Address     string            `json:"address"`
	Weight      int               `json:"weight"`
	Metadata    map[string]string `json:"metadata"`
	StartTime   int64             `json:"start_time"`
}

type ServiceEventType string

const (
	EventRegister   ServiceEventType = "register"
	EventDeregister ServiceEventType = "deregister"
)

type ServiceEvent struct {
	Type      ServiceEventType `json:"type"`
	Service   ServiceInfo      `json:"service"`
	Timestamp int64            `json:"timestamp"`
}

type client struct {
	cfg        Config
	httpClient *http.Client
	authToken  string
	authExpire time.Time
	authMu     sync.Mutex
}

func newClient(cfg Config) *client {
	if cfg.Group == "" {
		cfg.Group = "DEFAULT_GROUP"
	}
	if cfg.APIVersion == "" {
		cfg.APIVersion = "v3"
	}
	return &client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *client) namingEndpoint() string {
	if c.cfg.NamingEndpoint != "" {
		return c.cfg.NamingEndpoint
	}
	return c.cfg.Endpoint
}

func (c *client) ensureToken(ctx context.Context) (string, error) {
	if c.cfg.Username == "" || c.cfg.Password == "" {
		return c.cfg.Token, nil
	}
	c.authMu.Lock()
	defer c.authMu.Unlock()
	if c.authToken != "" && time.Now().Before(c.authExpire.Add(-60*time.Second)) {
		return c.authToken, nil
	}
	loginURL := fmt.Sprintf("%s/v1/auth/users/login", c.cfg.Endpoint)
	form := fmt.Sprintf("username=%s&password=%s", c.cfg.Username, c.cfg.Password)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, strings.NewReader(form))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("nacos login status %d", resp.StatusCode)
	}
	authHeader := resp.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		c.authToken = strings.TrimPrefix(authHeader, "Bearer ")
		c.authExpire = time.Now().Add(18000 * time.Second)
		return c.authToken, nil
	}
	body, _ := io.ReadAll(resp.Body)
	var loginResp struct {
		AccessToken string `json:"accessToken"`
	}
	if err := json.Unmarshal(body, &loginResp); err == nil && loginResp.AccessToken != "" {
		c.authToken = loginResp.AccessToken
		c.authExpire = time.Now().Add(18000 * time.Second)
		return c.authToken, nil
	}
	return "", fmt.Errorf("nacos login: no token in response")
}

func (c *client) doRequest(ctx context.Context, method, urlStr string, body io.Reader) (*http.Response, error) {
	token, _ := c.ensureToken(ctx)
	req, err := http.NewRequestWithContext(ctx, method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return c.httpClient.Do(req)
}

func parseAddrPort(addr string) (string, string) {
	idx := strings.LastIndex(addr, ":")
	if idx < 0 {
		return addr, "0"
	}
	return addr[:idx], addr[idx+1:]
}

func pollInterval(d string) time.Duration {
	if d == "" {
		return 5 * time.Second
	}
_dur, err := time.ParseDuration(d)
	if err != nil {
		return 5 * time.Second
	}
	if _dur < time.Second {
		return time.Second
	}
	return _dur
}

func maskURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if u.User != nil {
		u.User = url.User("****")
	}
	return u.String()
}
