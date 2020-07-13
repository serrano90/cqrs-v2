package memory_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/serrano90/cqrs-v2"
	"github.com/serrano90/cqrs-v2/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewInstanceOfDispatcher(t *testing.T) {
	var d cqrs.Dispatcher
	d = memory.NewDispatcherInMemory()
	if d == nil {
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
				NewTestCommand(),
			},
			expected: nil,
		},
		"when the type name exist": {
			handler: NewMockCommandHandler(),
			commandQuery: []interface{}{
				NewTestCommand(),
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
		expected     error
	}{
		"success when the values is a command": {
			handler:      NewMockCommandHandler(),
			commandQuery: NewTestCommand(),
			expected:     nil,
		},
		"success when the values is a query": {
			handler:      NewMockQueryHandler(),
			commandQuery: NewTestQuery(),
			expected:     nil,
		},
		"when the type name does not exit": {
			handler:      nil,
			commandQuery: NewTestCommand(),
			expected:     errors.New(cqrs.ErrMessageHandlerDoesNotExist),
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

		_, err := d.Dispatch(test.commandQuery)
		assert.Equal(t, test.expected, err, "They value does not equals")
	}
}

func NewTestCommand() cqrs.Command {
	return &TestCommand{}
}

type TestCommand struct {
	Id string
}

func (tc *TestCommand) TypeOf() string {
	return reflect.TypeOf(tc).Name()
}

func NewMockCommandHandler() cqrs.CommandHandler {
	return &MockCommandHandler{}
}

type MockCommandHandler struct{}

func (handle *MockCommandHandler) Handle(c cqrs.Command) (interface{}, error) {
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

func (handle *MockQueryHandler) Handle(c cqrs.Query) (interface{}, error) {
	return nil, nil
}
