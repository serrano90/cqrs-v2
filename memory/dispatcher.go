package memory

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/serrano90/cqrs-v2"
)

type DispatcherInMemory struct {
	middlewares []cqrs.CommandHandlerMiddleware
	handlers    map[string]interface{}
}

func NewDispatcherInMemory() cqrs.Dispatcher {
	return &DispatcherInMemory{
		middlewares: make([]cqrs.CommandHandlerMiddleware, 0),
		handlers:    make(map[string]interface{}, 0),
	}
}

func (d *DispatcherInMemory) Dispatch(ctx context.Context, cq interface{}) (interface{}, error) {
	typeOfName := d.getTypeOf(cq)
	if handler, ok := d.handlers[typeOfName]; ok {
		switch handler.(type) {
		case cqrs.CommandHandler:
			return d.dispatchToCommandHandler(ctx, cq, handler)
		case cqrs.QueryHandler:
			return d.dispatchToQueryHandler(ctx, cq, handler)
		}
	}
	return nil, errors.New(cqrs.ErrMessageHandlerDoesNotExist)
}

func (d *DispatcherInMemory) dispatchToCommandHandler(ctx context.Context, cq interface{}, handler interface{}) (interface{}, error) {
	ch, _ := handler.(cqrs.CommandHandler)
	h := ch.Handle
	c, _ := cq.(cqrs.Command)
	for _, m := range d.middlewares {
		h = m(h)
	}
	return h(ctx, c)
}

func (d *DispatcherInMemory) dispatchToQueryHandler(ctx context.Context, cq interface{}, handler interface{}) (interface{}, error) {
	h, _ := handler.(cqrs.QueryHandler)
	q, _ := cq.(cqrs.Command)
	return h.Handle(ctx, q)
}

func (d *DispatcherInMemory) AddHandler(handler interface{}, cq ...interface{}) error {
	for _, item := range cq {
		typeName := d.getTypeOf(item)
		if _, ok := d.handlers[typeName]; ok {
			return fmt.Errorf("%s %s", cqrs.ErrMessageHandlerDuplicated, typeName)
		}
		d.handlers[typeName] = handler
	}
	return nil
}

func (d *DispatcherInMemory) getTypeOf(cq interface{}) string {
	return reflect.TypeOf(cq).Elem().Name()
}

func (d *DispatcherInMemory) Use(middleware ...cqrs.CommandHandlerMiddleware) {
	d.middlewares = append(d.middlewares, middleware...)
}
