package gevent

import (
	"fmt"
	"reflect"
	"runtime"
)

type method struct {
	receiver reflect.Value
	method   reflect.Value
	typ      reflect.Type
	codec    CodecFunc
}

func newMethod(receiver reflect.Value, function interface{}, codec CodecFunc) (*method, error) {
	funcValue := reflect.ValueOf(function)
	if reflect.Func != funcValue.Kind() {
		return nil, fmt.Errorf("type (%v) is not func type", funcValue.Kind())
	}
	return &method{receiver: receiver, method: funcValue, typ: reflect.TypeOf(function), codec: codec}, nil
}

func (m *method) call(event Event, in ...interface{}) ([]reflect.Value, error) {
	var values []reflect.Value
	if m.receiver.IsValid() {
		values = append(values, m.receiver)
	}

	if m.codec != nil {
		var err error
		in, err = m.codec(in...)
		if err != nil {
			return nil, err
		}
	}

	if len(in)+len(values) != m.typ.NumIn() {
		return nil, fmt.Errorf("arguments nums error, len (%v) in call to [%v %#v]",
			len(in)+len(values), runtime.FuncForPC(m.method.Pointer()).Name(), m.method)
	}

	index := len(values)
	for i := 0; i < len(in) && index < m.typ.NumIn(); i++ {
		mt := m.typ.In(index)
		if in[i] == nil {
			values = append(values, reflect.New(mt).Elem())
		} else {
			if reflect.TypeOf(in[i]) != mt {
				if mt.Kind() != reflect.Interface ||
					!reflect.TypeOf(in[i]).Implements(mt) {
					return nil, fmt.Errorf("cannot use type (%v) as (%v) in argument to [%v %#v] args[%v]",
						reflect.TypeOf(in[i]), mt, runtime.FuncForPC(m.method.Pointer()).Name(), m, i)
				}
			}
			values = append(values, reflect.ValueOf(in[i]))
		}
		index++
	}

	outs := m.method.Call(values)
	return outs, nil
}
