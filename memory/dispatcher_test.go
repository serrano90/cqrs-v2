package memory_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/serrano90/cqrs-v2"
	"github.com/serrano90/cqrs-v2/memory"
	"github.com/serrano90/cqrs-v2/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNewInstanceOfDispatcher(t *testing.T) {
	if d := memory.NewDispatcherInMemory(); d == nil {
		t.Fail()
	}
}

func TestDispatcherAddHandler(t *testing.T) {
	tests := map[string]struct {
		handler      interface{}
		commandQuery []interface{}
		expected     error
	}{
		"success": {
			handler: NewMockCommandHandler(),
			commandQuery: []interface{}{
				NewTestCommand("x"),
			},
			expected: nil,
		},
		"when the type name exist": {
			handler: NewMockCommandHandler(),
			commandQuery: []interface{}{
				NewTestCommand("x"),
			},
			expected: errors.New(cqrs.ErrMessageHandlerDuplicated + " TestCommand"),
		},
	}

	d := memory.NewDispatcherInMemory()
	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		err := d.AddHandler(test.handler, test.commandQuery...)
		assert.Equal(t, test.expected, err, "They value does not equals")
	}

}

func TestDispatcherDispatch(t *testing.T) {
	tests := map[string]struct {
		handler      interface{}
		commandQuery interface{}
		middleware   []cqrs.CommandHandlerMiddleware
		expected     error
	}{
		"success when the values is a command": {
			handler:      NewMockCommandHandler(),
			commandQuery: NewTestCommand("x"),
			middleware:   nil,
			expected:     nil,
		},
		"success using middlewares": {
			handler:      NewMockCommandHandler(),
			commandQuery: NewTestCommand("x"),
			middleware: []cqrs.CommandHandlerMiddleware{
				middleware.NewValidationMiddleware(),
			},
			expected: nil,
		},
		"success when the values is a query": {
			handler:      NewMockQueryHandler(),
			commandQuery: NewTestQuery(),
			middleware:   nil,
			expected:     nil,
		},
		"when the type name does not exit": {
			handler:      nil,
			commandQuery: NewTestCommand("x"),
			middleware:   nil,
			expected:     errors.New(cqrs.ErrMessageHandlerDoesNotExist),
		},
		"when using middlewares and value is not valid": {
			handler:      NewMockCommandHandler(),
			commandQuery: NewTestCommand(""),
			middleware: []cqrs.CommandHandlerMiddleware{
				middleware.NewValidationMiddleware(),
			},
			expected: errors.New("The value is empty"),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		d := memory.NewDispatcherInMemory()
		if test.handler != nil {
			err := d.AddHandler(test.handler, test.commandQuery)
			if err != nil {
				t.Fail()
			}
		}

		d.Use(test.middleware...)

		_, err := d.Dispatch(context.Background(), test.commandQuery)
		assert.Equal(t, test.expected, err, "They value does not equals")
	}
}

func NewTestCommand(id string) cqrs.Command {
	return &TestCommand{
		Id: id,
	}
}

type TestCommand struct {
	Id string
}

func (tc *TestCommand) TypeOf() string {
	return reflect.TypeOf(tc).Name()
}

func (tc *TestCommand) Validate() error {
	if tc.Id == "" {
		return errors.New("The value is empty")
	}
	return nil
}

func NewMockCommandHandler() cqrs.CommandHandler {
	return &MockCommandHandler{}
}

type MockCommandHandler struct{}

func (handle *MockCommandHandler) Handle(ctx context.Context, c cqrs.Command) (interface{}, error) {
	return nil, nil
}

func NewTestQuery() cqrs.Command {
	return &TestQuery{}
}

type TestQuery struct {
	Id string
}

func (tc *TestQuery) TypeOf() string {
	return reflect.TypeOf(tc).Name()
}

func NewMockQueryHandler() cqrs.QueryHandler {
	return &MockQueryHandler{}
}

type MockQueryHandler struct{}

func (handle *MockQueryHandler) Handle(ctx context.Context, c cqrs.Query) (interface{}, error) {
	return nil, nil
}
