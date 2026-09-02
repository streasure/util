package timeutil

import "time"

const (
	SECONDS_OF_DAY = 86400
)

var _, localOffset = time.Now().Zone()
var offset64 = int64(localOffset)

var offset time.Duration

func SetOffset(d time.Duration) time.Duration {
	offset = d
	return offset
}

func GetOffset() time.Duration {
	return offset
}

func Now() time.Time {
	return time.Now().Add(offset)
}

func Timestamp() int64 {
	return Now().Unix()
}

func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}

func TimeStampToString(ts int64) string {
	tm := time.Unix(ts, 0)
	return tm.Format("2006-01-02 15:04:05")
}

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func DiffNatureDays(t1, t2 int64) int {
	if t1 == t2 {
		return 0
	}
	d1 := (t1 + offset64) / SECONDS_OF_DAY
	d2 := (t2 + offset64) / SECONDS_OF_DAY
	diff := d2 - d1
	if diff < 0 {
		return int(-diff)
	}
	return int(diff)
}

func DiffDays(endTime, startTime time.Time) int {
	start := startTime.Truncate(24 * time.Hour)
	end := endTime.Truncate(24 * time.Hour)
	return int(end.Sub(start).Hours() / 24)
}

func ZeroTimeOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func NormalizeTimeOfDay(t time.Time, startHour int) time.Time {
	y, m, d := t.Date()
	hour, _, _ := t.Clock()
	if hour < startHour {
		return time.Date(y, m, d-1, startHour, 0, 0, 0, t.Location())
	}
	return time.Date(y, m, d, startHour, 0, 0, 0, t.Location())
}

func GetTomorrowStamp() time.Time {
	now := Now()
	y, m, d := now.Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, now.Location())
}

func IsSameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func IsSameDayUnix(t1, t2 int64) bool {
	d1 := (t1 + offset64) / SECONDS_OF_DAY
	d2 := (t2 + offset64) / SECONDS_OF_DAY
	return d1 == d2
}

func IsSameWeek(t1, t2 int64) bool {
	y1, w1 := time.Unix(t1, 0).ISOWeek()
	y2, w2 := time.Unix(t2, 0).ISOWeek()
	return y1 == y2 && w1 == w2
}

func IsSameMonth(t1, t2 int64) bool {
	tm1 := time.Unix(t1, 0)
	tm2 := time.Unix(t2, 0)
	return tm1.Year() == tm2.Year() && tm1.Month() == tm2.Month()
}

func IsToday(t time.Time) bool {
	return IsSameDay(t, Now())
}

func IsTodayUnix(unix int64) bool {
	return IsSameDayUnix(unix, Timestamp())
}

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func GetTimeByHour(d time.Time, hour int) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), hour, 0, 0, 0, d.Location())
}

func GetDateKey(t time.Time) string {
	return t.Format("060102")
}
