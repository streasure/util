package errors

import (
	"testing"
)

func TestNew(t *testing.T) {
	err := New("test error")
	if err.Error() != "test error" {
		t.Errorf("New error = %q, want %q", err.Error(), "test error")
	}
}

func TestNewWithStack(t *testing.T) {
	err := NewWithStack("stack error")
	if err.Error() != "stack error" {
		t.Errorf("NewWithStack error = %q", err.Error())
	}
}

func TestWithMessage(t *testing.T) {
	err := New("original")
	wrapped := WithMessage(err, "context")
	if wrapped.Error() != "context: original" {
		t.Errorf("WithMessage error = %q", wrapped.Error())
	}
}

func TestWithStack(t *testing.T) {
	err := New("test")
	stacked := WithStack(err)
	if stacked.Error() != "test" {
		t.Errorf("WithStack error = %q", stacked.Error())
	}
}

func TestCode(t *testing.T) {
	err := NewECode(42)
	if Code(err) != 42 {
		t.Errorf("Code = %d, want 42", Code(err))
	}
	if Code(nil) != 0 {
		t.Errorf("Code(nil) = %d, want 0", Code(nil))
	}
}

func TestIs(t *testing.T) {
	err := New("target")
	wrapped := WithMessage(err, "wrap")
	if !Is(wrapped, err) {
		t.Error("Is should match wrapped error")
	}
}

func TestErrorsString(t *testing.T) {
	err1 := New("err1")
	err2 := New("err2")
	result := ErrorsString(err1, err2)
	if result == "" {
		t.Error("ErrorsString returned empty")
	}
}
