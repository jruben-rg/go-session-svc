package command

import (
	"context"
	"fmt"
	"testing"

	"github.com/jruben-rg/go-session-svc/sessions/domain/session"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestDeleteRepository struct {
	session.Repository
	err     error
	val     int64
	invoked bool
}

func (tdr *TestDeleteRepository) Delete(ctx context.Context, key string) (int64, error) {
	tdr.invoked = true
	return tdr.val, tdr.err
}

type TestDeleteSession deleteSessionHandler

func TestDeleteSessionHandlerShouldInvokeDeleteMethod(t *testing.T) {
	t.Parallel()

	logger := logrus.NewEntry(logrus.StandardLogger())

	tests := []struct {
		scenario        string
		expectedErr     error
		isErrorExpected bool
	}{
		{
			scenario:        "Should return error if repository returns error",
			expectedErr:     fmt.Errorf("Repository error"),
			isErrorExpected: true,
		},
		{
			scenario:        "Should not return error if repository does not return error",
			expectedErr:     nil,
			isErrorExpected: false,
		},
	}

	for _, test := range tests {

		repo := &TestDeleteRepository{err: test.expectedErr}
		handler := NewDeleteSessionHandler(repo, logger)
		err := handler.Handle(context.Background(), DeleteSession{Key: ""})

		if test.isErrorExpected {
			assert.NotNil(t, err, "An error is expected from the delete repository")
		} else {
			assert.Nil(t, err, "No error is expected from the delete repository")
		}

		assert.True(t, repo.invoked == true, "Delete method has been invoked")
	}

}

func TestDeleteSessionHandlerShouldPanicIfNilRepo(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Handle method did not panic")
		}
	}()

	logger := logrus.NewEntry(logrus.StandardLogger())
	handler := NewDeleteSessionHandler(nil, logger)
	handler.Handle(context.Background(), DeleteSession{Key: ""})

}
