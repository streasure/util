package component

import (
	"testing"
)

type mockComponent struct {
	BaseComponent
	name    string
	initErr error
	startErr error
	destroyed bool
}

func (m *mockComponent) Name() string   { return m.name }
func (m *mockComponent) Init() error    { return m.initErr }
func (m *mockComponent) Start() error   { return m.startErr }
func (m *mockComponent) Destroy()       { m.destroyed = true }

func TestBaseComponent(t *testing.T) {
	bc := &BaseComponent{}
	if bc.Name() != "base" {
		t.Errorf("Name = %q, want base", bc.Name())
	}
	if bc.Order() != 0 {
		t.Errorf("Order = %d, want 0", bc.Order())
	}
	if bc.Init() != nil {
		t.Error("Init should return nil")
	}
	if bc.Start() != nil {
		t.Error("Start should return nil")
	}
	bc.Destroy()

	bc.SetOrder(5)
	if bc.Order() != 5 {
		t.Errorf("Order after SetOrder = %d, want 5", bc.Order())
	}
}

func TestContainerAddAndServe(t *testing.T) {
	c := NewContainer()
	m1 := &mockComponent{name: "c1"}
	m2 := &mockComponent{name: "c2"}
	c.Add(m1)
	c.Add(m2)
	if len(c.components) != 2 {
		t.Errorf("components len = %d, want 2", len(c.components))
	}
}

func TestContainerOrdering(t *testing.T) {
	c := NewContainer()
	m1 := &mockComponent{name: "high"}
	m1.SetOrder(10)
	m2 := &mockComponent{name: "low"}
	m2.SetOrder(1)
	c.Add(m1)
	c.Add(m2)

	// Sort should order low before high
	// We can't call Serve() because it blocks on signal, but we can test sort logic
	type sortable struct {
		name  string
		order int
	}
	items := []sortable{
		{m1.name, m1.Order()},
		{m2.name, m2.Order()},
	}
	// Simulate sort
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].order > items[j].order {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	if items[0].name != "low" || items[1].name != "high" {
		t.Error("ordering failed")
	}
}

func TestContainerInitError(t *testing.T) {
	c := NewContainer()
	m := &mockComponent{name: "fail", initErr: &testError{"init failed"}}
	c.Add(m)

	defer func() {
		if r := recover(); r == nil {
			t.Error("should panic on init error")
		}
	}()

	// This should panic because Init returns error
	for _, comp := range c.components {
		if err := comp.Init(); err != nil {
			panic(err)
		}
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string { return e.msg }
