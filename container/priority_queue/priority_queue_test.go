package priority_queue

import "testing"

func TestMaxPriorityQueue(t *testing.T) {
	pq := New()

	pq.Push("low", 1)
	pq.Push("high", 10)
	pq.Push("medium", 5)

	if pq.Len() != 3 {
		t.Fatalf("Len: got %d, want 3", pq.Len())
	}

	v, p := pq.Pop()
	if v != "high" || p != 10 {
		t.Fatalf("Pop: got (%v, %d), want (high, 10)", v, p)
	}

	v, p = pq.Pop()
	if v != "medium" || p != 5 {
		t.Fatalf("Pop: got (%v, %d), want (medium, 5)", v, p)
	}

	v, p = pq.Pop()
	if v != "low" || p != 1 {
		t.Fatalf("Pop: got (%v, %d), want (low, 1)", v, p)
	}

	if pq.Len() != 0 {
		t.Fatalf("Len after popping all: got %d, want 0", pq.Len())
	}
}

func TestMinPriorityQueue(t *testing.T) {
	pq := New(WithMin(true))

	pq.Push("low", 1)
	pq.Push("high", 10)
	pq.Push("medium", 5)

	v, p := pq.Pop()
	if v != "low" || p != 1 {
		t.Fatalf("Pop: got (%v, %d), want (low, 1)", v, p)
	}

	v, p = pq.Pop()
	if v != "medium" || p != 5 {
		t.Fatalf("Pop: got (%v, %d), want (medium, 5)", v, p)
	}

	v, p = pq.Pop()
	if v != "high" || p != 10 {
		t.Fatalf("Pop: got (%v, %d), want (high, 10)", v, p)
	}
}

func TestPeek(t *testing.T) {
	pq := New()

	pq.Push("a", 1)
	pq.Push("b", 2)

	v, p := pq.Peek()
	if v != "b" || p != 2 {
		t.Fatalf("Peek: got (%v, %d), want (b, 2)", v, p)
	}

	if pq.Len() != 2 {
		t.Fatalf("Len after Peek: got %d, want 2", pq.Len())
	}
}

func TestPopEmpty(t *testing.T) {
	pq := New()

	v, p := pq.Pop()
	if v != nil || p != 0 {
		t.Fatalf("Pop empty: got (%v, %d), want (nil, 0)", v, p)
	}
}

func TestPeekEmpty(t *testing.T) {
	pq := New()

	v, p := pq.Peek()
	if v != nil || p != 0 {
		t.Fatalf("Peek empty: got (%v, %d), want (nil, 0)", v, p)
	}
}

func TestClear(t *testing.T) {
	pq := New()

	pq.Push("a", 1)
	pq.Push("b", 2)

	pq.Clear()
	if pq.Len() != 0 {
		t.Fatalf("Len after Clear: got %d, want 0", pq.Len())
	}
}

func TestManyElements(t *testing.T) {
	pq := New()

	for i := int64(100); i > 0; i-- {
		pq.Push(i, i)
	}

	if pq.Len() != 100 {
		t.Fatalf("Len: got %d, want 100", pq.Len())
	}

	for i := int64(100); i > 0; i-- {
		v, p := pq.Pop()
		if p != i {
			t.Fatalf("Pop: got priority %d, want %d", p, i)
		}
		if v.(int64) != i {
			t.Fatalf("Pop: got value %v, want %d", v, i)
		}
	}
}
