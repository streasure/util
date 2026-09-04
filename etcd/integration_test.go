package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdComponentIntegration(t *testing.T) {
	connection := Config{Endpoint: "http://127.0.0.1:2379", ServicePrefix: "/util-test/services", ConfigPrefix: "/util-test/config"}
	client, err := newClient(connection)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if _, err := client.Status(ctx, connection.endpointList()[0]); err != nil {
		t.Skipf("etcd is not available on 127.0.0.1:2379: %v", err)
	}

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	serviceID := "service-" + id
	configKey := "app-" + id + ".json"
	discovered := make(chan ServiceEvent, 2)
	updated := make(chan string, 4)
	discovery := New(ComponentConfig{
		Enabled:   true,
		Etcd:      connection,
		Discovery: DiscoveryConfig{Enabled: true, ServiceID: serviceID, LoadBalance: LoadBalanceRoundRobin},
		Config:    DynamicConfig{Enabled: true, Key: configKey, Format: "json"},
	})
	discovery.OnServiceChange(func(event ServiceEvent) { discovered <- event })
	discovery.OnConfigChange(func(data []byte) { updated <- string(data) })
	if err := discovery.Start(); err != nil {
		t.Fatal(err)
	}
	defer discovery.Destroy()

	registration := New(ComponentConfig{
		Enabled:      true,
		Etcd:         connection,
		Registration: RegistrationConfig{Enabled: true, ServiceID: serviceID, InstanceID: "instance-1", Address: "127.0.0.1:8080", LeaseTTL: "3s"},
	})
	if err := registration.Start(); err != nil {
		t.Fatal(err)
	}
	defer registration.Destroy()

	select {
	case event := <-discovered:
		if event.Type != EventRegister || event.InstanceID != "instance-1" || event.Address != "127.0.0.1:8080" {
			t.Fatalf("unexpected service event: %#v", event)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for service discovery")
	}
	if address, ok := discovery.Pick(); !ok || address != "127.0.0.1:8080" {
		t.Fatalf("unexpected selected service: %q, %v", address, ok)
	}

	leaseCtx, leaseCancel := context.WithTimeout(context.Background(), time.Second)
	registration.registrationMu.RLock()
	leaseID := registration.leaseID
	registration.registrationMu.RUnlock()
	lease, err := client.TimeToLive(leaseCtx, leaseID, clientv3.WithAttachedKeys())
	leaseCancel()
	if err != nil || lease.TTL <= 0 || len(lease.Keys) != 1 {
		t.Fatalf("invalid registration lease: %#v, %v", lease, err)
	}
	revokeCtx, revokeCancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Revoke(revokeCtx, leaseID)
	revokeCancel()
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for {
		registration.registrationMu.RLock()
		newLeaseID := registration.leaseID
		registration.registrationMu.RUnlock()
		getCtx, getCancel := context.WithTimeout(context.Background(), time.Second)
		response, getErr := client.Get(getCtx, registration.registrationKey())
		getCancel()
		if getErr == nil && newLeaseID != leaseID && len(response.Kvs) == 1 && response.Kvs[0].Lease == int64(newLeaseID) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("service was not re-registered after lease loss")
		}
		time.Sleep(100 * time.Millisecond)
	}

	putCtx, putCancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(putCtx, discovery.configKey(), `{"name":"enabled"}`)
	putCancel()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), time.Second)
		defer cleanupCancel()
		_, _ = client.Delete(cleanupCtx, discovery.configKey())
	}()
	select {
	case data := <-updated:
		if data != `{"name":"enabled"}` {
			t.Fatalf("unexpected dynamic config: %q", data)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for dynamic config")
	}

	registration.Destroy()
	select {
	case event := <-discovered:
		if event.Type != EventDeregister || event.InstanceID != "instance-1" {
			t.Fatalf("unexpected deregistration event: %#v", event)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for service deregistration")
	}
}
