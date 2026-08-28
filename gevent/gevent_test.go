package gevent

import (
	"testing"
)

type testService struct {
	value int
}

func (s *testService) GetValue() int {
	return s.value
}

func (s *testService) SetValue(v int) {
	s.value = v
}

func TestDispatcherRegisterAndDispatch(t *testing.T) {
	d := NewDispatcher()
	called := false

	err := d.Register("test.event", func() {
		called = true
	})
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}

	err = d.Dispatch("test.event")
	if err != nil {
		t.Fatalf("Dispatch error: %v", err)
	}

	if !called {
		t.Error("Handler was not called")
	}
}

func TestDispatcherCall(t *testing.T) {
	d := NewDispatcher()

	err := d.Register("test.add", func(a, b int) int {
		return a + b
	})
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}

	values, err := d.Call("test.add", 3, 4)
	if err != nil {
		t.Fatalf("Call error: %v", err)
	}

	if values[0].Int() != 7 {
		t.Errorf("Call result = %d, want 7", values[0].Int())
	}
}

func TestDispatcherNoEvent(t *testing.T) {
	d := NewDispatcher()
	err := d.Dispatch("nonexistent")
	if err == nil {
		t.Error("Dispatch nonexistent should error")
	}

	d2 := NewDispatcher(WithSuccessWhenNoEvent(true))
	err = d2.Dispatch("nonexistent")
	if err != nil {
		t.Errorf("WithSuccessWhenNoEvent should not error: %v", err)
	}
}

func TestDispatcherRegisterService(t *testing.T) {
	d := NewDispatcher()
	svc := &testService{value: 42}

	err := d.RegisterService(svc)
	if err != nil {
		t.Fatalf("RegisterService error: %v", err)
	}

	events := d.Dump()
	if len(events) == 0 {
		t.Error("Dump should return registered events")
	}
}

func TestNoEventError(t *testing.T) {
	err := NoEventError{Event: "test"}
	if !IsNoEventError(err) {
		t.Error("IsNoEventError should return true")
	}
}
