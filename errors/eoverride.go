package errors

import (
	"fmt"
	"io"
)

type eOverride struct {
	new error
	old error
}

func (e eOverride) Unwrap() error {
	return e.new
}

func (e eOverride) Error() string {
	return e.new.Error()
}

func (e eOverride) Cause() error {
	return e.new
}

func (e eOverride) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", e.new)
			if e.old != nil {
				fmt.Fprintf(s, "override:%+v", e.old)
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}
