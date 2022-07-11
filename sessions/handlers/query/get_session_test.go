package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/jruben-rg/go-session-svc/sessions/domain/session"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestGetRepository struct {
	session.Repository
	err     error
	value   interface{}
	invoked bool
}

func (tgr *TestGetRepository) Get(ctx context.Context, key string) (interface{}, error) {
	tgr.invoked = true
	return tgr.value, tgr.err
}

type TestGetSession getSessionHandler

func TestGetSessionHandlerShouldInvokeGetMethod(t *testing.T) {
	t.Parallel()

	logger := logrus.NewEntry(logrus.StandardLogger())

	tests := []struct {
		scenario        string
		expectedErr     error
		expectedVal     interface{}
		isErrorExpected bool
	}{
		{
			scenario:        "Should return error if repository returns error",
			expectedErr:     fmt.Errorf("Repository error"),
			expectedVal:     nil,
			isErrorExpected: true,
		},
		{
			scenario:        "Should not return error if repository does not return error",
			expectedErr:     nil,
			expectedVal:     `{"expected":"value"}`,
			isErrorExpected: false,
		},
	}

	for _, test := range tests {

		repo := &TestGetRepository{value: test.expectedVal, err: test.expectedErr}
		handler := NewGetSessionHandler(repo, logger)
		val, err := handler.Handle(context.Background(), GetSession{})

		if test.isErrorExpected {
			assert.NotNil(t, err, "An error is expected from the Get repository")
		} else {
			assert.Nil(t, err, "No error is expected from the Get repository")
		}

		assert.True(t, val == test.expectedVal, "Value from Get method matches expected result")
		assert.True(t, repo.invoked == true, "Get method has been invoked")
	}

}

func TestSetSessionHandlerShouldPanicIfNilRepo(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Handle method did not panic")
		}
	}()

	logger := logrus.NewEntry(logrus.StandardLogger())
	handler := NewGetSessionHandler(nil, logger)
	handler.Handle(context.Background(), GetSession{})

}
