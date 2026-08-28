package util

import "time"

const (
	SECONDS_OF_DAY = 86400
)

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
	if t1 > t2 {
		t1, t2 = t2, t1
	}

	diffDays := 0
	secDiff := t2 - t1
	if secDiff > SECONDS_OF_DAY {
		tmpDays := int(secDiff / SECONDS_OF_DAY)
		t1 += int64(tmpDays) * SECONDS_OF_DAY
		diffDays += tmpDays
	}

	st := time.Unix(t1, 0)
	et := time.Unix(t2, 0)
	dateFormatTpl := "20060102"
	if st.Format(dateFormatTpl) != et.Format(dateFormatTpl) {
		diffDays += 1
	}

	return diffDays
}

func IsSameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func IsSameDayUnix(t1, t2 int64) bool {
	tt1 := time.Unix(t1, 0)
	tt2 := time.Unix(t2, 0)
	return tt1.Year() == tt2.Year() && tt1.YearDay() == tt2.YearDay()
}

func IsSameWeek(t1, t2 int64) bool {
	y1, w1 := time.Unix(t1, 0).ISOWeek()
	y2, w2 := time.Unix(t2, 0).ISOWeek()
	return y1 == y2 && w1 == w2
}

func IsSameMonth(t1, t2 int64) bool {
	y1, m1, _ := time.Unix(t1, 0).Date()
	y2, m2, _ := time.Unix(t2, 0).Date()
	return y1 == y2 && m1 == m2
}

func IsToday(t time.Time) bool {
	return IsSameDay(t, time.Now())
}

func IsTodayUnix(unix int64) bool {
	return IsSameDay(time.Unix(unix, 0), time.Now())
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
