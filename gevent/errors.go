package gevent

import (
	"errors"
	"fmt"
)

type NoEventError struct {
	Event Event
}

func (e NoEventError) Error() string {
	return fmt.Sprintf("event(%v) not found", e.Event)
}

func IsNoEventError(err error) bool {
	return errors.As(err, &NoEventError{})
}
