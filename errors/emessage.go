package errors

import (
	"fmt"
	"io"
)

type eMessage struct {
	raw error
	msg string
}

func (e eMessage) Unwrap() error {
	return e.raw
}

func (e eMessage) Error() string {
	return e.msg + ": " + e.raw.Error()
}

func (e eMessage) Cause() error {
	return e.raw
}

func (e eMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s: ", e.msg)
			fmt.Fprintf(s, "%+v", e.Cause())
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}
