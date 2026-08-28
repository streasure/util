package gametime

import "time"

const (
	absoluteZeroYear = -292277022399
	daysPer400Years  = 365*400 + 97
	daysPer100Years  = 365*100 + 24
	daysPer4Years    = 365*4 + 1
)

var defaultTime = NewRefTime(DailyTime{0, 0, 0})

type DailyTime struct {
	Hour int
	Min  int
	Sec  int
}

type RefTime struct {
	DailyTime
	value int
}

func NewRefTime(dt DailyTime) *RefTime {
	return &RefTime{
		DailyTime: dt,
		value:     clockValue(dt.Hour, dt.Min, dt.Sec),
	}
}

func NewRefTimeDuration(d time.Duration) *RefTime {
	d = d % Day
	hour := int(d.Hours())
	min := int(d.Minutes()) % 60
	sec := int(d.Seconds()) % 60
	return NewRefTime(DailyTime{Hour: hour, Min: min, Sec: sec})
}

func NextNDayResetTime(t any, days int) time.Time {
	return defaultTime.NextNDayResetTime(t, days)
}

func NextNWeeksResetTime(t any, weeks int) time.Time {
	return defaultTime.NextNWeeksResetTime(t, weeks)
}

func IsSameDay(a, b any) bool {
	return defaultTime.IsSameDay(a, b)
}

func IsToday(t any) bool {
	return defaultTime.IsSameDay(t, Now())
}

func IsSameWeek(a, b any) bool {
	return defaultTime.IsSameWeek(a, b)
}

func IsSameWeekday(a, b any, weekday time.Weekday) bool {
	return defaultTime.IsSameWeekday(a, b, weekday)
}

func IsSameMonth(a, b any) bool {
	return defaultTime.IsSameMonth(a, b)
}

func IsSameMonthday(a, b any, monthday int) bool {
	return defaultTime.IsSameMonthday(a, b, monthday)
}

func IsCurrentWeek(t any) bool {
	return defaultTime.IsSameWeek(t, Now())
}

func SubDay(a, b any) int {
	return defaultTime.SubDay(a, b)
}

func (rt *RefTime) NextNDayResetTime(t any, days int) time.Time {
	tm := rt.parse(t)
	year, month, day := tm.Date()
	return time.Date(year, month, day-rt.dayOffset(tm)+days, rt.Hour, rt.Min, rt.Sec, 0, time.Local)
}

func (rt *RefTime) NextNWeeksResetTime(t any, weeks int) time.Time {
	return rt.NextNWeeksWeekdayResetTime(t, weeks, time.Monday)
}

func (rt *RefTime) NextNWeeksWeekdayResetTime(t any, weeks int, weekday time.Weekday) time.Time {
	tm := rt.parse(t)
	dayOffset := int(weekday-tm.Weekday()) + rt.dayOffset(tm)
	if dayOffset > 0 {
		dayOffset -= 7
	}
	return rt.NextNDayResetTime(tm, dayOffset+weeks*7)
}

func (rt *RefTime) NextNMonthsMonthdayResetTime(t any, months int, monthday int) time.Time {
	tm := rt.parse(t)
	year, month, day := tm.Date()

	var monthOffset int
	if day-rt.dayOffset(tm) < monthday {
		monthOffset = -1
	}

	return time.Date(year, month+time.Month(months+monthOffset), monthday, rt.Hour, rt.Min, rt.Sec, 0, time.Local)
}

func (rt *RefTime) IsToday(t any) bool {
	return rt.IsSameDay(t, Now())
}

func (rt *RefTime) IsCurrentWeek(t any) bool {
	return rt.IsSameWeek(t, Now())
}

func (rt *RefTime) IsSameDay(a, b any) bool {
	ta, tb := parseT(a), parseT(b)
	return rt.daysSinceEpoch(ta) == rt.daysSinceEpoch(tb)
}

func (rt *RefTime) IsSameWeek(a, b any) bool {
	return rt.IsSameWeekday(a, b, time.Monday)
}

func (rt *RefTime) IsSameWeekday(a, b any, weekday time.Weekday) bool {
	ta, tb, equal := rt.parseEqual(a, b)
	if equal {
		return true
	}
	return rt.NextNWeeksWeekdayResetTime(ta, 0, weekday).Equal(
		rt.NextNWeeksWeekdayResetTime(tb, 0, weekday))
}

func (rt *RefTime) IsSameMonth(a, b any) bool {
	return rt.IsSameMonthday(a, b, 1)
}

func (rt *RefTime) IsSameMonthday(a, b any, monthday int) bool {
	ta, tb, equal := rt.parseEqual(a, b)
	if equal {
		return true
	}
	return rt.NextNMonthsMonthdayResetTime(ta, 0, monthday).Equal(
		rt.NextNMonthsMonthdayResetTime(tb, 0, monthday))
}

func (rt *RefTime) SubDay(a, b any) int {
	return rt.daysSinceEpoch(rt.parse(a)) - rt.daysSinceEpoch(rt.parse(b))
}

func (rt *RefTime) dayOffset(t time.Time) int {
	if rt.value == 0 {
		return 0
	}
	hour, min, sec := t.Clock()
	if clockValue(hour, min, sec) < rt.value {
		return 1
	}
	return 0
}

func (rt *RefTime) parse(t any) time.Time {
	tm := parseT(t)
	if tm.Location() != time.Local {
		tm = tm.In(time.Local)
	}
	return tm
}

func (rt *RefTime) daysSinceEpoch(t time.Time) int {
	return daysSinceEpoch(t) - rt.dayOffset(t)
}

func (rt *RefTime) parseEqual(a, b any) (time.Time, time.Time, bool) {
	at, bt := rt.parse(a), rt.parse(b)
	return at, bt, at.Equal(bt)
}

func clockValue(hour, min, sec int) int {
	return hour*10000 + min*100 + sec
}

func daysSinceEpoch(t time.Time) int {
	y := t.Year() - absoluteZeroYear
	n := y / 400
	y -= 400 * n
	d := daysPer400Years * n

	n = y / 100
	y -= 100 * n
	d += daysPer100Years * n

	n = y / 4
	y -= 4 * n
	d += daysPer4Years * n

	n = y
	d += 365 * n

	return d + t.YearDay()
}

func parseT(t any) time.Time {
	switch v := t.(type) {
	case time.Time:
		return v
	case int64:
		return time.Unix(v, 0)
	case int32:
		return time.Unix(int64(v), 0)
	case int:
		return time.Unix(int64(v), 0)
	default:
		return time.Time{}
	}
}
