package now

import (
	"testing"
	"time"

	gametime "github.com/streasure/util/gametime"
)

func TestBeginningOfDay(t *testing.T) {
	gametime.SetOffset(0)
	now := gametime.Now()
	result := BeginningOfDay()
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Error("BeginningOfDay should return midnight")
	}
	if result.Day() != now.Day() {
		t.Error("BeginningOfDay should be same day")
	}
}

func TestBeginningOfWeek(t *testing.T) {
	gametime.SetOffset(0)
	result := BeginningOfWeek()
	if result.Weekday() != weekStartDay {
		t.Errorf("BeginningOfWeek should start on %v, got %v", weekStartDay, result.Weekday())
	}
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Error("BeginningOfWeek should return start of day")
	}
}

func TestBeginningOfMonth(t *testing.T) {
	gametime.SetOffset(0)
	result := BeginningOfMonth()
	if result.Day() != 1 {
		t.Errorf("BeginningOfMonth should be day 1, got %d", result.Day())
	}
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Error("BeginningOfMonth should return start of day")
	}
}

func TestBeginningOfYear(t *testing.T) {
	gametime.SetOffset(0)
	result := BeginningOfYear()
	if result.Month() != time.January || result.Day() != 1 {
		t.Errorf("BeginningOfYear should be Jan 1, got %v", result)
	}
	if result.Hour() != 0 || result.Minute() != 0 || result.Second() != 0 {
		t.Error("BeginningOfYear should return start of day")
	}
}

func TestEndOfDay(t *testing.T) {
	gametime.SetOffset(0)
	result := EndOfDay()
	if result.Hour() != 23 || result.Minute() != 59 || result.Second() != 59 {
		t.Error("EndOfDay should return 23:59:59")
	}
}

func TestEndOfWeek(t *testing.T) {
	gametime.SetOffset(0)
	result := EndOfWeek()
	if result.Hour() != 23 || result.Minute() != 59 || result.Second() != 59 {
		t.Error("EndOfWeek should return end of day")
	}
}

func TestEndOfMonth(t *testing.T) {
	gametime.SetOffset(0)
	result := EndOfMonth()
	if result.Hour() != 23 || result.Minute() != 59 || result.Second() != 59 {
		t.Error("EndOfMonth should return end of day")
	}
}

func TestEndOfYear(t *testing.T) {
	gametime.SetOffset(0)
	result := EndOfYear()
	if result.Month() != time.December || result.Day() != 31 {
		t.Errorf("EndOfYear should be Dec 31, got %v", result)
	}
	if result.Hour() != 23 || result.Minute() != 59 || result.Second() != 59 {
		t.Error("EndOfYear should return end of day")
	}
}

func TestWith(t *testing.T) {
	tm := time.Date(2021, 1, 15, 12, 30, 0, 0, time.Local)
	result := With(tm).BeginningOfDay()
	if result.Day() != 15 {
		t.Errorf("With().BeginningOfDay() day = %d, want 15", result.Day())
	}
}

func TestNew(t *testing.T) {
	tm := time.Date(2021, 1, 15, 12, 30, 0, 0, time.Local)
	result := New(tm).EndOfDay()
	if result.Day() != 15 {
		t.Errorf("New().EndOfDay() day = %d, want 15", result.Day())
	}
	if result.Hour() != 23 {
		t.Errorf("New().EndOfDay() hour = %d, want 23", result.Hour())
	}
}

func TestSetWeekStartDay(t *testing.T) {
	orig := GetWeekStartDay()
	SetWeekStartDay(time.Sunday)
	if GetWeekStartDay() != time.Sunday {
		t.Error("SetWeekStartDay failed")
	}
	SetWeekStartDay(orig)
}

func TestBeginningOfWeek_Sunday(t *testing.T) {
	SetWeekStartDay(time.Sunday)
	result := BeginningOfWeek()
	if result.Weekday() != time.Sunday {
		t.Errorf("BeginningOfWeek(Sunday) should start on Sunday, got %v", result.Weekday())
	}
}
