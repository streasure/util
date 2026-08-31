package wrand

import (
	"math/rand"
	"testing"
)

func TestPick(t *testing.T) {
	chooser, err := NewRandChooser(
		NewRandChoice("a", 1),
		NewRandChoice("b", 2),
		NewRandChoice("c", 3),
	)
	if err != nil {
		t.Fatalf("NewRandChooser error: %v", err)
	}

	counts := map[string]int{}
	for i := 0; i < 10000; i++ {
		item := chooser.Pick()
		counts[item]++
	}

	if counts["c"] <= counts["b"] {
		t.Fatal("c should be picked more often than b")
	}
	if counts["b"] <= counts["a"] {
		t.Fatal("b should be picked more often than a")
	}
}

func TestPickSource(t *testing.T) {
	chooser, err := NewRandChooser(
		NewRandChoice("a", 1),
		NewRandChoice("b", 1),
	)
	if err != nil {
		t.Fatalf("NewRandChooser error: %v", err)
	}

	rs := rand.New(rand.NewSource(42))
	counts := map[string]int{}
	for i := 0; i < 10000; i++ {
		item := chooser.PickSource(rs)
		counts[item]++
	}

	if counts["a"] == 0 || counts["b"] == 0 {
		t.Fatal("both items should be picked")
	}
}

func TestPickN(t *testing.T) {
	chooser, err := NewRandChooser(
		NewRandChoice("a", 1),
		NewRandChoice("b", 1),
		NewRandChoice("c", 1),
	)
	if err != nil {
		t.Fatalf("NewRandChooser error: %v", err)
	}

	picked := chooser.PickN(2)
	if len(picked) != 2 {
		t.Fatalf("PickN: got %d items, want 2", len(picked))
	}

	if picked[0] == picked[1] {
		t.Fatal("PickN should return different items")
	}
}

func TestPickNMoreThanAvailable(t *testing.T) {
	chooser, err := NewRandChooser(
		NewRandChoice("a", 1),
		NewRandChoice("b", 1),
	)
	if err != nil {
		t.Fatalf("NewRandChooser error: %v", err)
	}

	picked := chooser.PickN(10)
	if len(picked) != 2 {
		t.Fatalf("PickN over: got %d items, want 2", len(picked))
	}
}

func TestPickNSingle(t *testing.T) {
	chooser, err := NewRandChooser(
		NewRandChoice("a", 1),
		NewRandChoice("b", 1),
	)
	if err != nil {
		t.Fatalf("NewRandChooser error: %v", err)
	}

	picked := chooser.PickN(1)
	if len(picked) != 1 {
		t.Fatalf("PickN(1): got %d items, want 1", len(picked))
	}
}

func TestPickNZeroOrNegative(t *testing.T) {
	chooser, err := NewRandChooser(
		NewRandChoice("a", 1),
	)
	if err != nil {
		t.Fatalf("NewRandChooser error: %v", err)
	}

	picked := chooser.PickN(0)
	if picked != nil {
		t.Fatal("PickN(0) should return nil")
	}

	picked = chooser.PickN(-1)
	if picked != nil {
		t.Fatal("PickN(-1) should return nil")
	}
}

func TestNewRandChoices(t *testing.T) {
	choices := NewRandChoices([][]int{
		{1, 10},
		{2, 20},
		{3, 30},
	})

	if len(choices) != 3 {
		t.Fatalf("NewRandChoices: got %d, want 3", len(choices))
	}

	if choices[0].Item != 1 || choices[0].Weight != 10 {
		t.Fatalf("choices[0]: got (%d, %d), want (1, 10)", choices[0].Item, choices[0].Weight)
	}
}

func TestNewRandChoicesEmpty(t *testing.T) {
	choices := NewRandChoices([][]int{})
	if choices != nil {
		t.Fatal("NewRandChoices empty should return nil")
	}
}

func TestNewRandChoicesInvalid(t *testing.T) {
	choices := NewRandChoices([][]int{
		{1},
		{2, 3},
	})
	if len(choices) != 1 {
		t.Fatalf("NewRandChoices invalid: got %d, want 1", len(choices))
	}
}

func TestErrWeightOverflow(t *testing.T) {
	_, err := NewRandChooser(
		NewRandChoice("a", int(^uint(0)>>1)),
		NewRandChoice("b", int(^uint(0)>>1)),
	)
	if err != errWeightOverflow {
		t.Fatalf("expected errWeightOverflow, got: %v", err)
	}
}

func TestErrNoValidChoices(t *testing.T) {
	_, err := NewRandChooser(
		NewRandChoice("a", 0),
		NewRandChoice("b", 0),
	)
	if err != errNoValidChoices {
		t.Fatalf("expected errNoValidChoices, got: %v", err)
	}
}
