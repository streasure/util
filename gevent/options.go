package gevent

import (
	"fmt"
	"reflect"
)

type Options struct {
	DefaultRegisterOptions RegisterOptions
	ErrorHandle            errorHandle
	RegisterErrorPanic     bool
	SuccessWhenNoEvent     bool
}

type Option func(o *Options)

func WithDefaultRegisterOptions(opts ...RegisterOption) Option {
	registerOptions := defaultRegisterOptions()
	for _, o := range opts {
		o(&registerOptions)
	}
	return func(o *Options) {
		o.DefaultRegisterOptions = registerOptions
	}
}

type errorHandle func(error)

func WithErrorHandle(h errorHandle) Option {
	return func(o *Options) {
		o.ErrorHandle = h
	}
}

func WithRegisterErrorPanic(b bool) Option {
	return func(o *Options) {
		o.RegisterErrorPanic = b
	}
}

func WithSuccessWhenNoEvent(b bool) Option {
	return func(o *Options) {
		o.SuccessWhenNoEvent = b
	}
}

type RegisterOptions struct {
	Unique         bool
	EventGenerator EventGeneratorFunc
	FilterMethod   FilterMethodFunc
	CodecGenerator CodecGeneratorFunc
	Codec          CodecFunc
}

func defaultRegisterOptions() RegisterOptions {
	return RegisterOptions{
		EventGenerator: DefaultEventGenerator,
	}
}

type RegisterOption func(o *RegisterOptions)

func RegisterWithUnique(b bool) RegisterOption {
	return func(o *RegisterOptions) {
		o.Unique = b
	}
}

var DefaultEventGenerator = func(receiver interface{}, serviceName string, method reflect.Method) Event {
	return fmt.Sprintf("%s.%s", serviceName, method.Name)
}

type EventGeneratorFunc func(receiver interface{}, serviceName string, method reflect.Method) Event

func RegisterWithEventGenerator(f EventGeneratorFunc) RegisterOption {
	return func(o *RegisterOptions) {
		o.EventGenerator = f
	}
}

type FilterMethodFunc func(methodType reflect.Type) bool

func RegisterWithFilterMethod(f FilterMethodFunc) RegisterOption {
	return func(o *RegisterOptions) {
		o.FilterMethod = f
	}
}

type CodecGeneratorFunc func(receiver interface{}, serviceName string, method reflect.Method) CodecFunc
type CodecFunc func(arguments ...interface{}) ([]interface{}, error)

func RegisterWithCodecGenerator(f CodecGeneratorFunc) RegisterOption {
	return func(o *RegisterOptions) {
		o.CodecGenerator = f
	}
}
