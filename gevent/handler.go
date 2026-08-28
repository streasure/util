package gevent

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Event interface{}

type handler struct {
	event   Event
	methods []*method
}

func newHandler(event Event) *handler {
	return &handler{event: event}
}

func (h *handler) trigger(arguments ...interface{}) error {
	if len(h.methods) == 0 {
		return NoEventError{h.event}
	}

	var errs []string
	for _, m := range h.methods {
		if _, err := m.call(h.event, arguments...); err != nil {
			errs = append(errs, fmt.Errorf("trigger event(%v) error: %w", h.event, err).Error())
		}
	}

	if len(errs) != 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func (h *handler) call(arguments ...interface{}) ([]reflect.Value, error) {
	if len(h.methods) == 0 {
		return nil, NoEventError{h.event}
	}
	if len(h.methods) > 1 {
		return nil, fmt.Errorf("event(%v) has multi handler", h.event)
	}
	if values, err := h.methods[0].call(h.event, arguments...); err != nil {
		return nil, fmt.Errorf("call event(%v) error: %w", h.event, err)
	} else {
		return values, nil
	}
}

func (h *handler) register(receiver reflect.Value, handlerFunc interface{}, opts RegisterOptions) error {
	if opts.Unique && len(h.methods) > 0 {
		return fmt.Errorf("register event(%v) error: event has been registered", h.event)
	}
	method, err := newMethod(receiver, handlerFunc, opts.Codec)
	if err != nil {
		return fmt.Errorf("register event(%v) error: %w", h.event, err)
	}
	h.methods = append(h.methods, method)
	return nil
}
