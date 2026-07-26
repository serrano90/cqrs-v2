package memory_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/serrano90/cqrs-v2/v3"
	"github.com/serrano90/cqrs-v2/v3/memory"
	"github.com/serrano90/cqrs-v2/v3/middleware"
	"github.com/stretchr/testify/assert"
)

func TestNewInstanceOfDispatcher(t *testing.T) {
	if d := memory.NewDispatcherInMemory(); d == nil {
		t.Fail()
	}
}

func TestDispatcherAddCommandHandler(t *testing.T) {
	d := memory.NewDispatcherInMemory()

	// First registration should succeed
	err := memory.AddCommandHandler[*TestCommand, any](d, NewMockCommandHandler())
	assert.Equal(t, nil, err, "First registration should succeed")

	// Second registration for the same command type should return duplicated error
	err = memory.AddCommandHandler[*TestCommand, any](d, NewMockCommandHandler())
	assert.Equal(t, errors.New(cqrs.ErrMessageHandlerDuplicated+" TestCommand"), err, "Should return duplicated error")
}

func TestDispatcherDispatchCommand(t *testing.T) {
	tests := map[string]struct {
		handler    cqrs.CommandHandler[*TestCommand, any]
		command    *TestCommand
		middleware []cqrs.CommandHandlerMiddleware[*TestCommand, any]
		expected   error
	}{
		"success when the value is a command": {
			handler:    NewMockCommandHandler(),
			command:    NewTestCommand("x"),
			middleware: nil,
			expected:   nil,
		},
		"success using middlewares": {
			handler: NewMockCommandHandler(),
			command: NewTestCommand("x"),
			middleware: []cqrs.CommandHandlerMiddleware[*TestCommand, any]{
				middleware.NewValidationMiddleware[*TestCommand, any](),
			},
			expected: nil,
		},
		"when the handler does not exist": {
			handler:    nil,
			command:    NewTestCommand("x"),
			middleware: nil,
			expected:   errors.New(cqrs.ErrMessageHandlerDoesNotExist),
		},
		"when using middlewares and value is not valid": {
			handler: NewMockCommandHandler(),
			command: NewTestCommand(""),
			middleware: []cqrs.CommandHandlerMiddleware[*TestCommand, any]{
				middleware.NewValidationMiddleware[*TestCommand, any](),
			},
			expected: errors.New("The value is empty"),
		},
	}

	for name, test := range tests {
		t.Logf("Running test case: %s", name)

		d := memory.NewDispatcherInMemory()
		if test.handler != nil {
			err := memory.AddCommandHandler(d, test.handler)
			if err != nil {
				t.Fail()
			}
		}

		for _, m := range test.middleware {
			memory.Use(d, m)
		}

		_, err := memory.DispatchCommand[*TestCommand, any](context.Background(), d, test.command)
		assert.Equal(t, test.expected, err, "The value does not equal")
	}
}

func TestDispatcherDispatchQuery(t *testing.T) {
	d := memory.NewDispatcherInMemory()

	err := memory.AddQueryHandler[*TestQuery, string](d, NewMockQueryHandler())
	assert.Equal(t, nil, err, "Registration should succeed")

	res, err := memory.DispatchQuery[*TestQuery, string](context.Background(), d, NewTestQuery())
	assert.Equal(t, nil, err, "Dispatch should succeed")
	assert.Equal(t, "result", res, "Should return the handler result")

	// Dispatching a query without a registered handler should fail
	d2 := memory.NewDispatcherInMemory()
	_, err = memory.DispatchQuery[*TestQuery, string](context.Background(), d2, NewTestQuery())
	assert.Equal(t, errors.New(cqrs.ErrMessageHandlerDoesNotExist), err, "Should return does-not-exist error")
}

func TestDispatcherDispatchIncompatibleResultType(t *testing.T) {
	d := memory.NewDispatcherInMemory()

	err := memory.AddCommandHandler[*TestCommand, any](d, NewMockCommandHandler())
	assert.Equal(t, nil, err, "Registration should succeed")

	// Dispatching with a result type different from the registered one fails
	_, err = memory.DispatchCommand[*TestCommand, string](context.Background(), d, NewTestCommand("x"))
	assert.Equal(t, errors.New("registered handler has incompatible type for TestCommand"), err, "Should return incompatible-type error")
}

func NewTestCommand(id string) *TestCommand {
	return &TestCommand{
		ID: id,
	}
}

type TestCommand struct {
	ID string
}

func (tc *TestCommand) TypeOf() string {
	return reflect.TypeOf(tc).Elem().Name()
}

func (tc *TestCommand) Validate() error {
	if tc.ID == "" {
		return errors.New("The value is empty")
	}
	return nil
}

func NewMockCommandHandler() cqrs.CommandHandler[*TestCommand, any] {
	return &MockCommandHandler{}
}

type MockCommandHandler struct{}

func (handle *MockCommandHandler) Handle(_ context.Context, _ *TestCommand) (any, error) {
	return nil, nil
}

func NewTestQuery() *TestQuery {
	return &TestQuery{}
}

type TestQuery struct {
	ID string
}

func (tc *TestQuery) TypeOf() string {
	return reflect.TypeOf(tc).Elem().Name()
}

func NewMockQueryHandler() cqrs.QueryHandler[*TestQuery, string] {
	return &MockQueryHandler{}
}

type MockQueryHandler struct{}

func (handle *MockQueryHandler) Handle(_ context.Context, _ *TestQuery) (string, error) {
	return "result", nil
}
