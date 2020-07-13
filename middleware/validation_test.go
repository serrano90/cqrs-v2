package middleware_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/serrano90/cqrs-v2"
	"github.com/serrano90/cqrs-v2/middleware"
	"github.com/stretchr/testify/assert"
)

type TestCommand struct {
	value string
}

func (tc *TestCommand) TypeOf() string {
	return reflect.TypeOf(tc).Name()
}

func (tc *TestCommand) Validate() error {
	if tc.value == "" {
		return errors.New("The value is empty")
	}
	return nil
}

type TestCommandHandler struct{}

func (h *TestCommandHandler) Handle(cmd cqrs.Command) (interface{}, error) {
	return "Success", nil
}

func TestValidateMiddleware(t *testing.T) {
	cases := map[string]struct {
		cmd           cqrs.Command
		expected      interface{}
		expectedError error
	}{
		"success": {
			cmd: &TestCommand{
				value: "value",
			},
			expected:      "Success",
			expectedError: nil,
		},
		"when the value does not valid": {
			cmd: &TestCommand{
				value: "",
			},
			expected:      nil,
			expectedError: errors.New("The value is empty"),
		},
	}

	ch := &TestCommandHandler{}
	m := middleware.NewValidationMiddleware()
	for name, test := range cases {
		t.Logf("Running test case: %s", name)
		h := m(ch.Handle)
		resp, err := h(test.cmd)

		assert.Equal(t, test.expected, resp, "The expected value and result value not equals")
		assert.Equal(t, test.expectedError, err, "The expected value and result value not equals")
	}
}
