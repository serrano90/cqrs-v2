package middleware_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/serrano90/cqrs-v2/v3/middleware"
	"github.com/stretchr/testify/assert"
)

type TestCommand struct {
	value string
}

func (tc *TestCommand) TypeOf() string {
	return reflect.TypeOf(tc).Elem().Name()
}

func (tc *TestCommand) Validate() error {
	if tc.value == "" {
		return errors.New("The value is empty")
	}
	return nil
}

type TestCommandHandler struct{}

func (h *TestCommandHandler) Handle(_ context.Context, _ *TestCommand) (string, error) {
	return "Success", nil
}

func TestValidateMiddleware(t *testing.T) {
	cases := map[string]struct {
		cmd           *TestCommand
		expected      string
		expectedError error
	}{
		"success": {
			cmd: &TestCommand{
				value: "value",
			},
			expected:      "Success",
			expectedError: nil,
		},
		"when the value is not valid": {
			cmd: &TestCommand{
				value: "",
			},
			expected:      "",
			expectedError: errors.New("The value is empty"),
		},
	}

	ch := &TestCommandHandler{}
	m := middleware.NewValidationMiddleware[*TestCommand, string]()
	for name, test := range cases {
		t.Logf("Running test case: %s", name)
		h := m(ch.Handle)
		resp, err := h(context.Background(), test.cmd)

		assert.Equal(t, test.expected, resp, "The expected value and result value are not equal")
		assert.Equal(t, test.expectedError, err, "The expected error and result error are not equal")
	}
}
