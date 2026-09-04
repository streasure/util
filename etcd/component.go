package etcd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/streasure/util/component"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

type LoadBalance string

const (
	LoadBalanceRoundRobin LoadBalance = "round_robin"
	LoadBalanceRandom     LoadBalance = "random"
	LoadBalanceP2C        LoadBalance = "p2c"
)

type RegistrationConfig struct {
	Enabled    bool   `yaml:"enabled" json:"enabled"`
	ServiceID  string `yaml:"serviceID" json:"serviceID"`
	InstanceID string `yaml:"instanceID" json:"instanceID"`
	Address    string `yaml:"address" json:"address"`
	LeaseTTL   string `yaml:"leaseTTL" json:"leaseTTL"`
}

type DiscoveryConfig struct {
	Enabled     bool        `yaml:"enabled" json:"enabled"`
	ServiceID   string      `yaml:"serviceID" json:"serviceID"`
	LoadBalance LoadBalance `yaml:"loadBalance" json:"loadBalance"`
}

type DynamicConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Key     string `yaml:"key" json:"key"`
	Format  string `yaml:"format" json:"format"`
}

type ComponentConfig struct {
	Enabled      bool               `yaml:"enabled" json:"enabled"`
	Etcd         Config             `yaml:"etcd" json:"etcd"`
	Registration RegistrationConfig `yaml:"registration" json:"registration"`
	Discovery    DiscoveryConfig    `yaml:"discovery" json:"discovery"`
	Config       DynamicConfig      `yaml:"config" json:"config"`
}

type ServiceEventType string

const (
	EventRegister   ServiceEventType = "register"
	EventDeregister ServiceEventType = "deregister"
)

type ServiceEvent struct {
	Type       ServiceEventType `json:"type"`
	ServiceID  string           `json:"serviceID"`
	InstanceID string           `json:"instanceID"`
	Address    string           `json:"address"`
}

type ServiceChangeCallback func(ServiceEvent)
type ConfigChangeCallback func([]byte)
type ConfigValidator func([]byte) error
type ConfigValidate interface{ Validate() error }

func NewConfigValidator[T any](format string) ConfigValidator {
	return func(data []byte) error {
		var value T
		if err := DecodeConfig(data, format, &value); err != nil {
			return err
		}
		if validated, ok := any(&value).(ConfigValidate); ok {
			return validated.Validate()
		}
		return nil
	}
}

type Component struct {
	component.BaseComponent
	cfg              ComponentConfig
	client           *clientv3.Client
	serviceKey       string
	leaseID          clientv3.LeaseID
	registrationMu   sync.RWMutex
	services         atomic.Value // []string
	serviceSet       atomic.Value // map[string]string
	loads            sync.Map
	roundRobin       atomic.Uint64
	config           atomic.Value // []byte
	lastError        atomic.Value // string
	serviceCallbacks []ServiceChangeCallback
	configCallbacks  []ConfigChangeCallback
	validator        ConfigValidator
	mu               sync.RWMutex
	lifecycleMu      sync.Mutex
	started          bool
	ctx              context.Context
	cancel           context.CancelFunc
	wg               sync.WaitGroup
}

func New(cfg ComponentConfig) *Component {
	return &Component{cfg: cfg}
}

func (c *Component) Name() string { return "etcd" }
func (c *Component) Init() error  { return c.validate() }

func (c *Component) Start() error {
	c.lifecycleMu.Lock()
	defer c.lifecycleMu.Unlock()
	if c.started {
		return fmt.Errorf("etcd component already started")
	}
	c.started = true
	if !c.cfg.Enabled {
		return nil
	}
	if err := c.validate(); err != nil {
		c.started = false
		return err
	}
	client, err := newClient(c.cfg.Etcd)
	if err != nil {
		c.started = false
		return err
	}
	c.client = client
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.services.Store([]string{})
	c.serviceSet.Store(map[string]string{})
	c.config.Store([]byte{})
	c.lastError.Store("")

	if c.cfg.Registration.Enabled {
		if err := c.register(); err != nil {
			return c.failStart(err)
		}
		c.wg.Add(1)
		go c.keepAliveLoop()
	}
	if c.cfg.Discovery.Enabled {
		if err := c.refreshServices(); err != nil {
			return c.failStart(err)
		}
		c.wg.Add(1)
		go c.watchServices()
	}
	if c.cfg.Config.Enabled {
		if err := c.refreshConfig(); err != nil {
			return c.failStart(err)
		}
		c.wg.Add(1)
		go c.watchConfig()
	}
	return nil
}

func (c *Component) failStart(err error) error {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
	if c.client != nil {
		_ = c.client.Close()
	}
	c.started = false
	return err
}

func (c *Component) Destroy() {
	c.lifecycleMu.Lock()
	if !c.started {
		c.lifecycleMu.Unlock()
		return
	}
	c.started = false
	cancelComponent := c.cancel
	client := c.client
	c.lifecycleMu.Unlock()
	if cancelComponent != nil {
		cancelComponent()
	}
	c.wg.Wait()
	if client != nil {
		c.registrationMu.RLock()
		serviceKey := c.serviceKey
		c.registrationMu.RUnlock()
		if serviceKey != "" {
			ctx, cancel := requestContext()
			_, _ = client.Delete(ctx, serviceKey)
			cancel()
		}
		_ = client.Close()
	}
}

func (c *Component) OnServiceChange(callback ServiceChangeCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.serviceCallbacks = append(c.serviceCallbacks, callback)
}
func (c *Component) OnConfigChange(callback ConfigChangeCallback) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.configCallbacks = append(c.configCallbacks, callback)
}
func (c *Component) SetConfigValidator(validator ConfigValidator) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.validator = validator
}

func (c *Component) Services() []string {
	services, _ := c.services.Load().([]string)
	return append([]string(nil), services...)
}
func (c *Component) ConfigSnapshot() []byte {
	data, _ := c.config.Load().([]byte)
	return append([]byte(nil), data...)
}
func (c *Component) LastError() error {
	message, _ := c.lastError.Load().(string)
	if message == "" {
		return nil
	}
	return fmt.Errorf("%s", message)
}
func (c *Component) setError(err error) {
	if err == nil {
		c.lastError.Store("")
		return
	}
	c.lastError.Store(err.Error())
}

func DecodeConfig[T any](data []byte, format string, target *T) error {
	if target == nil {
		return fmt.Errorf("etcd config decode target is nil")
	}
	switch strings.ToLower(format) {
	case "json":
		return json.Unmarshal(data, target)
	case "yaml", "yml":
		return yaml.Unmarshal(data, target)
	case "toml":
		return toml.Unmarshal(data, target)
	default:
		return fmt.Errorf("unsupported config format %q", format)
	}
}

func BindConfig[T any](c *Component, target *T) error {
	if c == nil {
		return fmt.Errorf("etcd component is nil")
	}
	data := c.ConfigSnapshot()
	if err := DecodeConfig(data, c.cfg.Config.Format, target); err != nil {
		return err
	}
	if validated, ok := any(target).(ConfigValidate); ok {
		return validated.Validate()
	}
	return nil
}

func OnTypedConfigChange[T any](c *Component, callback func(*T) error) {
	if c == nil || callback == nil {
		return
	}
	c.SetConfigValidator(NewConfigValidator[T](c.cfg.Config.Format))
	c.OnConfigChange(func(data []byte) {
		var value T
		if err := DecodeConfig(data, c.cfg.Config.Format, &value); err == nil {
			if validated, ok := any(&value).(ConfigValidate); ok {
				if err := validated.Validate(); err != nil {
					return
				}
			}
			c.setError(callback(&value))
		}
	})
}

func (c *Component) Acquire() (string, func(), bool) {
	address, ok := c.Pick()
	if !ok {
		return "", nil, false
	}
	if c.cfg.Discovery.LoadBalance != LoadBalanceP2C {
		return address, func() {}, true
	}
	return address, c.Begin(address), true
}

func (c *Component) Pick() (string, bool) {
	services := c.Services()
	if len(services) == 0 {
		return "", false
	}
	switch c.cfg.Discovery.LoadBalance {
	case LoadBalanceRandom:
		return services[rand.IntN(len(services))], true
	case LoadBalanceP2C:
		if len(services) == 1 {
			return services[0], true
		}
		first := rand.IntN(len(services))
		second := rand.IntN(len(services) - 1)
		if second >= first {
			second++
		}
		if c.load(services[first]) <= c.load(services[second]) {
			return services[first], true
		}
		return services[second], true
	default:
		return services[(c.roundRobin.Add(1)-1)%uint64(len(services))], true
	}
}

func (c *Component) Begin(address string) func() {
	value, _ := c.loads.LoadOrStore(address, new(atomic.Int64))
	value.(*atomic.Int64).Add(1)
	var once sync.Once
	return func() { once.Do(func() { value.(*atomic.Int64).Add(-1) }) }
}
func (c *Component) load(address string) int64 {
	if value, ok := c.loads.Load(address); ok {
		return value.(*atomic.Int64).Load()
	}
	return 0
}

func (c *Component) validate() error {
	if !c.cfg.Enabled {
		return nil
	}
	if err := c.cfg.Etcd.validate(); err != nil {
		return err
	}
	if c.cfg.Registration.Enabled {
		if c.cfg.Registration.ServiceID == "" || c.cfg.Registration.InstanceID == "" || c.cfg.Registration.Address == "" {
			return fmt.Errorf("etcd registration serviceID, instanceID and address are required")
		}
		if _, err := leaseTTL(c.cfg.Registration.LeaseTTL); err != nil {
			return err
		}
	}
	if c.cfg.Discovery.Enabled && c.cfg.Discovery.ServiceID == "" {
		return fmt.Errorf("etcd discovery serviceID is required")
	}
	if c.cfg.Config.Enabled && c.cfg.Config.Key == "" {
		return fmt.Errorf("etcd dynamic config key is required")
	}
	if c.cfg.Config.Enabled {
		format := strings.ToLower(c.cfg.Config.Format)
		if format != "json" && format != "yaml" && format != "yml" && format != "toml" {
			return fmt.Errorf("unsupported dynamic config format %q", c.cfg.Config.Format)
		}
	}
	return nil
}
func leaseTTL(value string) (time.Duration, error) {
	if value == "" {
		return 10 * time.Second, nil
	}
	ttl, err := time.ParseDuration(value)
	if err != nil || ttl < time.Second {
		return 0, fmt.Errorf("etcd leaseTTL must be at least one second")
	}
	return ttl, nil
}
func (c *Component) registrationKey() string {
	return prefixJoin(c.cfg.Etcd.ServicePrefix, c.cfg.Registration.ServiceID) + "/" + c.cfg.Registration.InstanceID
}
func (c *Component) discoveryPrefix() string {
	return prefixJoin(c.cfg.Etcd.ServicePrefix, c.cfg.Discovery.ServiceID) + "/"
}
func (c *Component) configKey() string {
	if strings.HasPrefix(c.cfg.Config.Key, "/") {
		return c.cfg.Config.Key
	}
	return prefixJoin(c.cfg.Etcd.ConfigPrefix, c.cfg.Config.Key)
}

func (c *Component) register() error {
	ttl, err := leaseTTL(c.cfg.Registration.LeaseTTL)
	if err != nil {
		return err
	}
	ctx, cancel := requestContext()
	defer cancel()
	grant, err := c.client.Grant(ctx, int64(ttl/time.Second))
	if err != nil {
		return fmt.Errorf("etcd grant lease: %w", err)
	}
	key := c.registrationKey()
	if _, err = c.client.Put(ctx, key, c.cfg.Registration.Address, clientv3.WithLease(grant.ID)); err != nil {
		return fmt.Errorf("etcd register: %w", err)
	}
	c.registrationMu.Lock()
	c.leaseID, c.serviceKey = grant.ID, key
	c.registrationMu.Unlock()
	return nil
}
func (c *Component) keepAliveLoop() {
	defer c.wg.Done()
	for c.ctx.Err() == nil {
		c.registrationMu.RLock()
		leaseID := c.leaseID
		c.registrationMu.RUnlock()
		ch, err := c.client.KeepAlive(c.ctx, leaseID)
		if err == nil {
			for range ch {
				if c.ctx.Err() != nil {
					return
				}
			}
			if c.ctx.Err() != nil {
				return
			}
		}
		select {
		case <-c.ctx.Done():
			return
		case <-time.After(time.Second):
		}
		if err := c.register(); err != nil {
			c.setError(err)
			continue
		}
		c.setError(nil)
	}
}
func (c *Component) refreshServices() error {
	ctx, cancel := requestContext()
	defer cancel()
	resp, err := c.client.Get(ctx, c.discoveryPrefix(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	set := make(map[string]string, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		address := string(kv.Value)
		set[string(kv.Key)] = address
	}
	c.applyServiceSet(set)
	return nil
}
func (c *Component) watchServices() {
	defer c.wg.Done()
	revision := int64(0)
	for c.ctx.Err() == nil {
		ctx, cancel := requestContext()
		resp, err := c.client.Get(ctx, c.discoveryPrefix(), clientv3.WithPrefix())
		cancel()
		if err != nil {
			c.setError(err)
			if !c.retryWait() {
				return
			}
			continue
		}
		c.applyServiceSnapshot(resp.Kvs)
		revision = resp.Header.Revision + 1
		watch := c.client.Watch(c.ctx, c.discoveryPrefix(), clientv3.WithPrefix(), clientv3.WithRev(revision))
		for {
			var response clientv3.WatchResponse
			var ok bool
			select {
			case <-c.ctx.Done():
				return
			case response, ok = <-watch:
			}
			if !ok {
				break
			}
			if response.Err() != nil {
				break
			}
			c.applyServiceEvents(response.Events)
		}
		if c.ctx.Err() == nil && !c.retryWait() {
			return
		}
	}
}
func (c *Component) refreshConfig() error {
	ctx, cancel := requestContext()
	defer cancel()
	resp, err := c.client.Get(ctx, c.configKey())
	if err != nil {
		return err
	}
	var data []byte
	if len(resp.Kvs) > 0 {
		data = append([]byte(nil), resp.Kvs[0].Value...)
	}
	return c.storeConfig(data)
}
func (c *Component) watchConfig() {
	defer c.wg.Done()
	for c.ctx.Err() == nil {
		ctx, cancel := requestContext()
		resp, err := c.client.Get(ctx, c.configKey())
		cancel()
		if err != nil {
			c.setError(err)
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}
		var snapshot []byte
		if len(resp.Kvs) > 0 {
			snapshot = resp.Kvs[0].Value
		}
		if err := c.storeConfig(snapshot); err != nil {
			c.setError(err)
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(time.Second):
			}
			continue
		}
		watch := c.client.Watch(c.ctx, c.configKey(), clientv3.WithRev(resp.Header.Revision+1))
		for {
			var response clientv3.WatchResponse
			var ok bool
			select {
			case <-c.ctx.Done():
				return
			case response, ok = <-watch:
			}
			if !ok {
				break
			}
			if response.Err() != nil {
				break
			}
			for _, event := range response.Events {
				if event.Type == clientv3.EventTypeDelete {
					c.setError(c.storeConfig(nil))
				} else {
					c.setError(c.storeConfig(event.Kv.Value))
				}
			}
		}
		if c.ctx.Err() == nil {
			select {
			case <-c.ctx.Done():
				return
			case <-time.After(time.Second):
			}
		}
	}
}
func (c *Component) retryWait() bool {
	select {
	case <-c.ctx.Done():
		return false
	case <-time.After(time.Second):
		return true
	}
}
func (c *Component) storeConfig(data []byte) error {
	c.mu.RLock()
	validator := c.validator
	callbacks := append([]ConfigChangeCallback(nil), c.configCallbacks...)
	c.mu.RUnlock()
	if err := validateConfigData(data, c.cfg.Config.Format); err != nil {
		return err
	}
	if validator != nil {
		if err := validator(data); err != nil {
			return fmt.Errorf("validate etcd dynamic config: %w", err)
		}
	}
	previous, _ := c.config.Load().([]byte)
	if bytes.Equal(previous, data) {
		return nil
	}
	c.config.Store(append([]byte(nil), data...))
	for _, callback := range callbacks {
		func() { defer func() { _ = recover() }(); callback(append([]byte(nil), data...)) }()
	}
	return nil
}

func validateConfigData(data []byte, format string) error {
	if len(data) == 0 {
		return nil
	}
	var value any
	switch strings.ToLower(format) {
	case "json":
		if err := json.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("invalid JSON dynamic config: %w", err)
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("invalid YAML dynamic config: %w", err)
		}
	case "toml":
		if _, err := toml.Decode(string(data), &value); err != nil {
			return fmt.Errorf("invalid TOML dynamic config: %w", err)
		}
	}
	return nil
}
func (c *Component) applyServiceSnapshot(kvs []*mvccpb.KeyValue) {
	set := make(map[string]string, len(kvs))
	for _, kv := range kvs {
		set[string(kv.Key)] = string(kv.Value)
	}
	c.applyServiceSet(set)
}
func (c *Component) applyServiceEvents(events []*clientv3.Event) {
	current, _ := c.serviceSet.Load().(map[string]string)
	next := make(map[string]string, len(current))
	for key, value := range current {
		next[key] = value
	}
	for _, event := range events {
		key := string(event.Kv.Key)
		if event.Type == clientv3.EventTypeDelete {
			delete(next, key)
		} else {
			next[key] = string(event.Kv.Value)
		}
	}
	c.applyServiceSet(next)
}
func (c *Component) applyServiceSet(set map[string]string) {
	old, _ := c.serviceSet.Load().(map[string]string)
	services := make([]string, 0, len(set))
	for _, address := range set {
		services = append(services, address)
	}
	sort.Strings(services)
	c.serviceSet.Store(set)
	c.services.Store(services)
	c.mu.RLock()
	callbacks := append([]ServiceChangeCallback(nil), c.serviceCallbacks...)
	c.mu.RUnlock()
	for key, address := range set {
		oldAddress, existed := old[key]
		if !existed || oldAddress != address {
			if existed {
				c.notifyServiceCallbacks(callbacks, ServiceEvent{Type: EventDeregister, ServiceID: c.cfg.Discovery.ServiceID, InstanceID: key[strings.LastIndex(key, "/")+1:], Address: oldAddress})
			}
			for _, callback := range callbacks {
				func() {
					defer func() { _ = recover() }()
					callback(ServiceEvent{Type: EventRegister, ServiceID: c.cfg.Discovery.ServiceID, InstanceID: key[strings.LastIndex(key, "/")+1:], Address: address})
				}()
			}
		}
	}
	for key, address := range old {
		if _, ok := set[key]; !ok {
			c.loads.Delete(address)
			for _, callback := range callbacks {
				func() {
					defer func() { _ = recover() }()
					callback(ServiceEvent{Type: EventDeregister, ServiceID: c.cfg.Discovery.ServiceID, InstanceID: key[strings.LastIndex(key, "/")+1:], Address: address})
				}()
			}
		}
	}
}

func (c *Component) notifyServiceCallbacks(callbacks []ServiceChangeCallback, event ServiceEvent) {
	for _, callback := range callbacks {
		func() { defer func() { _ = recover() }(); callback(event) }()
	}
}
