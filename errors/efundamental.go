package errors

import (
	"fmt"
	"io"
)

type eFundamental struct {
	msg string
	*stack
}

func (e eFundamental) Error() string {
	return e.msg
}

func (e eFundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.msg)
			e.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, e.msg)
	case 'q':
		fmt.Fprintf(s, "%q", e.msg)
	}
}
