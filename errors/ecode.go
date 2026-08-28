package errors

import "fmt"

type ECoder interface {
	Code() int
}

func NewECode(code int) error {
	return eCode(code)
}

type eCode int

func (e eCode) Error() string {
	return fmt.Sprintf("[ECode: %d] ", e.Code())
}

func (e eCode) Code() int {
	return int(e)
}
