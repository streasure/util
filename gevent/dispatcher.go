package gevent

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Dispatcher interface {
	Register(event Event, handlerFunc interface{}, opts ...RegisterOption) error
	RegisterWithReceiver(event Event, receiver interface{}, handlerFunc interface{}, opts ...RegisterOption) error
	RegisterService(receiver interface{}, opts ...RegisterOption) error
	Dispatch(event Event, arguments ...interface{}) error
	Call(event Event, arguments ...interface{}) ([]reflect.Value, error)
	Dump() []Event
}

type dispatcher struct {
	sync.RWMutex
	handlers map[Event]*handler
	options  Options
}

func NewDispatcher(opts ...Option) Dispatcher {
	options := Options{
		DefaultRegisterOptions: defaultRegisterOptions(),
	}
	for _, o := range opts {
		o(&options)
	}
	return &dispatcher{
		handlers: make(map[Event]*handler),
		options:  options,
	}
}

func (d *dispatcher) Register(event Event, handlerFunc interface{}, opts ...RegisterOption) error {
	return d.RegisterWithReceiver(event, nil, handlerFunc, opts...)
}

func (d *dispatcher) RegisterWithReceiver(event Event, receiver interface{}, handlerFunc interface{}, opts ...RegisterOption) error {
	registerOptions := d.options.DefaultRegisterOptions
	for _, o := range opts {
		o(&registerOptions)
	}
	return d.register(event, receiver, handlerFunc, registerOptions)
}

func (d *dispatcher) RegisterService(receiver interface{}, opts ...RegisterOption) error {
	serviceName := reflect.Indirect(reflect.ValueOf(receiver)).Type().Name()
	return d.registerService(receiver, serviceName, opts...)
}

func (d *dispatcher) Dispatch(event Event, arguments ...interface{}) error {
	if hl, ok := d.handlers[event]; ok {
		if err := hl.trigger(arguments...); err != nil {
			return d.onError(err)
		}
		return nil
	}
	if d.options.SuccessWhenNoEvent {
		return nil
	}
	return d.onError(NoEventError{event})
}

func (d *dispatcher) Call(event Event, arguments ...interface{}) ([]reflect.Value, error) {
	if hl, ok := d.handlers[event]; ok {
		if values, err := hl.call(arguments...); err != nil {
			return nil, d.onError(err)
		} else {
			return values, nil
		}
	}
	return nil, d.onError(NoEventError{event})
}

func (d *dispatcher) Dump() []Event {
	events := make([]Event, 0, len(d.handlers))
	for k := range d.handlers {
		events = append(events, k)
	}
	return events
}

func (d *dispatcher) register(event Event, receiver interface{}, handlerFunc interface{}, opts RegisterOptions) error {
	rv := reflect.ValueOf(receiver)
	if receiver != nil {
		if rv.Kind() != reflect.Interface && rv.Kind() != reflect.Ptr {
			return fmt.Errorf("receiver must be nil or type must be reflect.Interface or reflect.Ptr, cannot use type (%v)", rv.Kind())
		}
		if rv.IsNil() {
			return fmt.Errorf("receiver(%v) Interface or Ptr IsNil", receiver)
		}
	}

	d.Lock()
	defer d.Unlock()

	if _, ok := d.handlers[event]; !ok {
		d.handlers[event] = newHandler(event)
	}

	if err := d.handlers[event].register(rv, handlerFunc, opts); err != nil {
		if d.options.RegisterErrorPanic {
			panic(d.onError(err))
		}
		return d.onError(err)
	}
	return nil
}

func (d *dispatcher) registerService(receiver interface{}, serviceName string, opts ...RegisterOption) error {
	registerOptions := d.options.DefaultRegisterOptions
	for _, o := range opts {
		o(&registerOptions)
	}

	typ := reflect.TypeOf(receiver)
	if serviceName == "" {
		return errors.New("gevent.Register: no service name for type " + typ.String())
	}

	if typ.NumMethod() == 0 {
		if reflect.PointerTo(typ).NumMethod() != 0 {
			return fmt.Errorf("gevent.Register service: type %s has no exported methods of suitable type"+
				"(hint: pass a pointer to value of that type)", serviceName)
		}
		return fmt.Errorf("gevent.Register service: type %s has no exported methods of suitable type", serviceName)
	}

	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if registerOptions.FilterMethod != nil && !registerOptions.FilterMethod(method.Type) {
			continue
		}
		if registerOptions.CodecGenerator != nil {
			registerOptions.Codec = registerOptions.CodecGenerator(receiver, serviceName, method)
		}
		event := registerOptions.EventGenerator(receiver, serviceName, method)
		err := d.register(event, receiver, method.Func.Interface(), registerOptions)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *dispatcher) onError(err error) error {
	if d.options.ErrorHandle != nil {
		d.options.ErrorHandle(err)
	}
	return err
}
