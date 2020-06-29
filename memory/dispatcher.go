package memory

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/serrano90/cqrs/v2"
)

type DispatcherInMemory struct {
	handlers map[string]interface{}
}

func NewDispatcherInMemory() cqrs.Dispatcher {
	return &DispatcherInMemory{
		handlers: make(map[string]interface{}, 0),
	}
}

func (d *DispatcherInMemory) Dispatch(cq interface{}) (interface{}, error) {
	typeOfName := d.getTypeOf(cq)
	if handler, ok := d.handlers[typeOfName]; ok {
		switch handler.(type) {
		case cqrs.CommandHandler:
			return d.dispatchToCommandHandler(cq, handler)
		case cqrs.QueryHandler:
			return d.dispatchToQueryHandler(cq, handler)
		}
	}
	return nil, errors.New(cqrs.ErrMessageHandlerDoesNotExist)
}

func (d *DispatcherInMemory) dispatchToCommandHandler(cq interface{}, handler interface{}) (interface{}, error) {
	h, _ := handler.(cqrs.CommandHandler)
	c, _ := cq.(cqrs.Command)
	return h.Handle(c)
}

func (d *DispatcherInMemory) dispatchToQueryHandler(cq interface{}, handler interface{}) (interface{}, error) {
	h, _ := handler.(cqrs.QueryHandler)
	q, _ := cq.(cqrs.Command)
	return h.Handle(q)
}

func (d *DispatcherInMemory) AddHandler(handler interface{}, cq ...interface{}) error {
	for _, item := range cq {
		typeName := d.getTypeOf(item)
		if _, ok := d.handlers[typeName]; ok {
			return errors.New(fmt.Sprintf("%s %s", cqrs.ErrMessageHandlerDuplicated, typeName))
		}
		d.handlers[typeName] = handler
	}
	return nil
}

func (d *DispatcherInMemory) getTypeOf(cq interface{}) string {
	return reflect.TypeOf(cq).Elem().Name()
}
