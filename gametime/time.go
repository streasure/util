package gametime

import "time"

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
	Nanosecond        = time.Nanosecond
	Microsecond       = time.Microsecond
	Millisecond       = time.Millisecond
	Second            = time.Second
	Minute            = time.Minute
	Hour              = time.Hour
	Day               = 24 * time.Hour
)

var (
	ZeroUtcTime = time.Unix(0, 0)
	offset      time.Duration
)

func GetOffset() time.Duration {
	return offset
}

func Unix(t int64, nsec int64) time.Time {
	return time.Unix(t, nsec)
}

func Timestamp() int64 {
	return Now().Unix()
}

func SetOffset(d time.Duration) time.Duration {
	offset = d
	return offset
}

func Now() time.Time {
	return time.Now().Add(offset)
}

func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}
