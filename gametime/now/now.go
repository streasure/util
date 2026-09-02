package now

import (
	"time"

	gametime "github.com/streasure/util/gametime"
)

var weekStartDay = time.Monday

func SetWeekStartDay(day time.Weekday) {
	weekStartDay = day
}

func GetWeekStartDay() time.Weekday {
	return weekStartDay
}

func With(t time.Time) *Now {
	return &Now{t: t}
}

func New(t time.Time) *Now {
	return &Now{t: t}
}

func BeginningOfDay() time.Time {
	return With(gametime.Now()).BeginningOfDay()
}

func BeginningOfWeek() time.Time {
	return With(gametime.Now()).BeginningOfWeek()
}

func BeginningOfMonth() time.Time {
	return With(gametime.Now()).BeginningOfMonth()
}

func BeginningOfYear() time.Time {
	return With(gametime.Now()).BeginningOfYear()
}

func EndOfDay() time.Time {
	return With(gametime.Now()).EndOfDay()
}

func EndOfWeek() time.Time {
	return With(gametime.Now()).EndOfWeek()
}

func EndOfMonth() time.Time {
	return With(gametime.Now()).EndOfMonth()
}

func EndOfYear() time.Time {
	return With(gametime.Now()).EndOfYear()
}

type Now struct {
	t time.Time
}

func (n *Now) BeginningOfDay() time.Time {
	y, m, d := n.t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, n.t.Location())
}

func (n *Now) BeginningOfWeek() time.Time {
	y, m, d := n.t.Date()
	weekday := n.t.Weekday()
	offset := int(weekday - weekStartDay)
	if offset < 0 {
		offset += 7
	}
	return time.Date(y, m, d-offset, 0, 0, 0, 0, n.t.Location())
}

func (n *Now) BeginningOfMonth() time.Time {
	y, m, _ := n.t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, n.t.Location())
}

func (n *Now) BeginningOfYear() time.Time {
	y, _, _ := n.t.Date()
	return time.Date(y, 1, 1, 0, 0, 0, 0, n.t.Location())
}

func (n *Now) EndOfDay() time.Time {
	y, m, d := n.t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), n.t.Location())
}

func (n *Now) EndOfWeek() time.Time {
	return n.BeginningOfWeek().AddDate(0, 0, 7).Add(-time.Nanosecond)
}

func (n *Now) EndOfMonth() time.Time {
	y, m, _ := n.t.Date()
	return time.Date(y, m+1, 0, 23, 59, 59, int(time.Second-time.Nanosecond), n.t.Location())
}

func (n *Now) EndOfYear() time.Time {
	y, _, _ := n.t.Date()
	return time.Date(y+1, 1, 0, 23, 59, 59, int(time.Second-time.Nanosecond), n.t.Location())
}

func (n *Now) Time() time.Time {
	return n.t
}
