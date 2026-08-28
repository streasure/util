package errors

import (
	"fmt"
	"io"
)

type eStack struct {
	raw   error
	stack *stack
}

func (e eStack) Unwrap() error {
	return e.raw
}

func (e eStack) Error() string {
	return e.raw.Error()
}

func (e eStack) Cause() error {
	return e.raw
}

func (e eStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%v", e.Cause())
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.Error())
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}
