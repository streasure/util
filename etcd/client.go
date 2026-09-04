package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Config struct {
	Endpoints     []string `yaml:"endpoints" json:"endpoints"`
	Endpoint      string   `yaml:"endpoint" json:"endpoint"`
	Username      string   `yaml:"username" json:"username"`
	Password      string   `yaml:"password" json:"password"`
	DialTimeout   string   `yaml:"dialTimeout" json:"dialTimeout"`
	ServicePrefix string   `yaml:"servicePrefix" json:"servicePrefix"`
	ConfigPrefix  string   `yaml:"configPrefix" json:"configPrefix"`
	CertFile      string   `yaml:"certFile" json:"certFile"`
	CertKeyFile   string   `yaml:"certKeyFile" json:"certKeyFile"`
	CACertFile    string   `yaml:"caCertFile" json:"caCertFile"`
	TLSServerName string   `yaml:"tlsServerName" json:"tlsServerName"`
}

func (c Config) validate() error {
	var scheme string
	for _, endpoint := range c.endpointList() {
		if endpoint == "" {
			return fmt.Errorf("etcd endpoint cannot be empty")
		}
		u, err := url.Parse(endpoint)
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
			return fmt.Errorf("invalid etcd endpoint %q", endpoint)
		}
		if scheme == "" {
			scheme = u.Scheme
		}
		if scheme != u.Scheme {
			return fmt.Errorf("etcd endpoints must use the same scheme")
		}
	}
	if (c.CertFile == "") != (c.CertKeyFile == "") {
		return fmt.Errorf("certFile and certKeyFile must be configured together")
	}
	if scheme == "http" && (c.CACertFile != "" || c.CertFile != "" || c.CertKeyFile != "" || c.TLSServerName != "") {
		return fmt.Errorf("TLS settings require https etcd endpoints")
	}
	return nil
}

func (c Config) endpointList() []string {
	if len(c.Endpoints) > 0 {
		endpoints := append([]string(nil), c.Endpoints...)
		for i := range endpoints {
			endpoints[i] = strings.TrimSpace(endpoints[i])
		}
		return endpoints
	}
	if c.Endpoint == "" {
		return []string{"http://127.0.0.1:2379"}
	}
	endpoints := strings.Split(c.Endpoint, ",")
	for i := range endpoints {
		endpoints[i] = strings.TrimSpace(endpoints[i])
	}
	return endpoints
}

func (c Config) timeout() time.Duration {
	if c.DialTimeout == "" {
		return 5 * time.Second
	}
	d, err := time.ParseDuration(c.DialTimeout)
	if err != nil || d <= 0 {
		return 5 * time.Second
	}
	return d
}

func (c Config) tlsConfig() (*tls.Config, error) {
	if c.CACertFile == "" && c.CertFile == "" && c.CertKeyFile == "" && c.TLSServerName == "" {
		return nil, nil
	}
	config := &tls.Config{MinVersion: tls.VersionTLS12, ServerName: c.TLSServerName}
	if (c.CertFile == "") != (c.CertKeyFile == "") {
		return nil, fmt.Errorf("both etcd certFile and certKeyFile are required for client certificate authentication")
	}
	if c.CertFile != "" {
		certificate, err := tls.LoadX509KeyPair(c.CertFile, c.CertKeyFile)
		if err != nil {
			return nil, fmt.Errorf("load etcd client certificate: %w", err)
		}
		config.Certificates = []tls.Certificate{certificate}
	}
	if c.CACertFile != "" {
		data, err := os.ReadFile(c.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("read etcd CA file: %w", err)
		}
		roots, err := x509.SystemCertPool()
		if err != nil {
			roots = x509.NewCertPool()
		}
		if !roots.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("etcd CA file contains no certificates")
		}
		config.RootCAs = roots
	}
	return config, nil
}

func newClient(cfg Config) (*clientv3.Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	tlsConfig, err := cfg.tlsConfig()
	if err != nil {
		return nil, err
	}
	return clientv3.New(clientv3.Config{
		Endpoints:        cfg.endpointList(),
		Username:         cfg.Username,
		Password:         cfg.Password,
		DialTimeout:      cfg.timeout(),
		TLS:              tlsConfig,
		AutoSyncInterval: 30 * time.Second,
	})
}

func requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func prefixJoin(prefix, value string) string {
	prefix = strings.TrimRight(prefix, "/")
	value = strings.Trim(value, "/")
	if prefix == "" {
		return "/" + value
	}
	return prefix + "/" + value
}
