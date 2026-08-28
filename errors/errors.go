package errors

import (
	"errors"
	"fmt"
	"strings"
)

const (
	successCode = 0
	unknownCode = -1
)

var (
	callerSkip = 3
	frameDepth = 1
)

func CallSkip(skip int) {
	callerSkip = skip
}

func FrameDepth(i int) {
	frameDepth = i
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func WithMessage(err error, msg string) error {
	if err == nil {
		return New(msg)
	}
	return eMessage{raw: err, msg: msg}
}

func WithMessagef(err error, format string, args ...any) error {
	return WithMessage(err, fmt.Sprintf(format, args...))
}

func New(msg string) error {
	return errors.New(msg)
}

func NewWithStack(msg string) error {
	return &eFundamental{
		msg:   msg,
		stack: callers(),
	}
}

func ErrorsString(errs ...error) string {
	if len(errs) == 0 {
		return ""
	}
	builder := strings.Builder{}
	builder.WriteString("multiErr:\n")
	for _, err := range errs {
		builder.WriteByte('\t')
		builder.WriteString(err.Error())
		builder.WriteString(";\n")
	}
	return builder.String()
}

func Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func ErrorfWithStack(format string, args ...any) error {
	return WithStack(fmt.Errorf(format, args...))
}

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if errors.As(err, &eStack{}) || errors.As(err, &eFundamental{}) {
		return err
	}
	return eStack{raw: err, stack: callers()}
}

func Append(newErr, oldErr error) error {
	switch {
	case newErr == nil && oldErr == nil:
		return nil
	case newErr == nil:
		return oldErr
	case oldErr == nil:
		return newErr
	default:
		return &eOverride{new: newErr, old: oldErr}
	}
}

func Wrap(newErr, oldErr error, message string) error {
	switch {
	case newErr == nil && oldErr == nil:
		return nil
	case newErr == nil:
		return WithMessage(oldErr, message)
	case oldErr == nil:
		return WithMessage(newErr, message)
	default:
		err := &eOverride{new: newErr, old: oldErr}
		return WithMessage(err, message)
	}
}

func Wrapf(newErr, oldErr error, format string, args ...interface{}) error {
	return Wrap(newErr, oldErr, fmt.Sprintf(format, args...))
}

func WithOverride(newErr, oldErr error) error {
	switch {
	case newErr == nil && oldErr == nil:
		return nil
	case newErr == nil:
		return oldErr
	case oldErr == nil:
		return newErr
	default:
		return &eOverride{new: newErr, old: oldErr}
	}
}

func Code(err error) int {
	if err == nil {
		return successCode
	}
	for {
		if err == nil {
			return unknownCode
		}
		e, ok := err.(ECoder)
		if ok {
			return e.Code()
		}
		err = Unwrap(err)
	}
}

func Cause(err error) error {
	type causer interface {
		Cause() error
	}
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}
