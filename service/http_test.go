package service_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jruben-rg/go-session-svc/service"
	"github.com/jruben-rg/go-session-svc/sessions/handlers"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/command"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/query"
	"github.com/stretchr/testify/assert"
)

type testExpectationsHttp struct {
	invoked    bool
	handlerVal interface{}
	handlerErr error
}

type DeleteSessionHandlerHttp struct {
	command.DeleteSessionHandler
	testExpectationsHttp
}

type SetSessionHandlerHttp struct {
	command.SetSessionHandler
	testExpectationsHttp
}

type GetSessionHandlerHttp struct {
	query.GetSessionHandler
	testExpectationsHttp
}

func (dsht *DeleteSessionHandlerHttp) Handle(ctx context.Context, cmd command.DeleteSession) error {
	dsht.invoked = true
	return dsht.handlerErr
}

func (ssht *SetSessionHandlerHttp) Handle(ctx context.Context, cmd command.SetSession) error {
	ssht.invoked = true
	return ssht.handlerErr
}

func (gsht *GetSessionHandlerHttp) Handle(ctx context.Context, cmd query.GetSession) (interface{}, error) {
	gsht.invoked = true
	return gsht.handlerVal, gsht.handlerErr
}

func TestSetHttpSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario        string
		expectedInvoked bool
		expectedStatus  int
		requestBody     io.Reader
		err             error
	}{
		{
			scenario:        "Should respond with bad request if request body empty",
			expectedInvoked: false,
			expectedStatus:  http.StatusBadRequest,
			requestBody:     nil,
			err:             nil,
		},
		{
			scenario:        "Should respond with bad request if session key empty",
			expectedInvoked: false,
			expectedStatus:  http.StatusBadRequest,
			requestBody:     strings.NewReader(`{"sessionKey":"","sessionValue":{"value":"test"}}`),
			err:             nil,
		},
		{
			scenario:        "Should respond with internal server error if handler returns an error",
			expectedInvoked: true,
			expectedStatus:  http.StatusInternalServerError,
			requestBody:     strings.NewReader(`{"sessionKey":"key","sessionValue":{"value":"test"}}`),
			err:             fmt.Errorf("Error from handler"),
		},
		{
			scenario:        "Should respond with accepted if no errors are found",
			expectedInvoked: true,
			expectedStatus:  http.StatusAccepted,
			requestBody:     strings.NewReader(`{"sessionKey":"key","sessionValue":{"value":"test"}}`),
			err:             nil,
		},
	}

	for _, test := range tests {

		setSessionHandler := &SetSessionHandlerGrpc{
			testExpectationsGrpc: testExpectationsGrpc{
				handlerErr: test.err,
			},
		}

		setCommand := handlers.Commands{
			SetSession: setSessionHandler,
		}

		appSet := handlers.Application{
			Commands: setCommand,
		}

		httpSvc := service.NewHttpService(appSet)

		request := httptest.NewRequest(http.MethodPost, "/api/session", test.requestBody)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		httpSvc.SetSession(response, request)

		assert.True(t, response.Code == test.expectedStatus, fmt.Sprintf("Should respond with status code %d\n", test.expectedStatus))

		if test.expectedInvoked {
			assert.True(t, setSessionHandler.invoked, "'Handle' should have been invoked")
		} else {
			assert.False(t, setSessionHandler.invoked, "'Handle' should not have been invoked")
		}
	}

}

func TestGetHttpSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario        string
		expectedInvoked bool
		expectedStatus  int
		sessionKey      string
		val             interface{}
		err             error
	}{
		{
			scenario:        "Should respond with bad request if SessionKey is empty",
			expectedInvoked: false,
			expectedStatus:  http.StatusBadRequest,
			sessionKey:      "",
			err:             nil,
		},
		{
			scenario:        "Should respond with internal server error if handler responds with error",
			expectedInvoked: true,
			expectedStatus:  http.StatusInternalServerError,
			sessionKey:      "sessionKeyValue",
			err:             fmt.Errorf("Error from handler"),
		},
		{
			scenario:        "Should respond with internal server error if result cannot be parsed",
			expectedInvoked: true,
			expectedStatus:  http.StatusInternalServerError,
			sessionKey:      "sessionKeyValue",
			val:             nil,
			err:             nil,
		},
		{
			scenario:        "Should respond with not found if result is empty",
			expectedInvoked: true,
			expectedStatus:  http.StatusNotFound,
			sessionKey:      "sessionKeyValue",
			val:             ``,
			err:             nil,
		},
		{
			scenario:        "Should respond with internal server error if handler response cannot be parsed",
			expectedInvoked: true,
			expectedStatus:  http.StatusInternalServerError,
			sessionKey:      "sessionKeyValue",
			val:             `aValue`,
			err:             nil,
		},
		{
			scenario:        "Should respond with session data",
			expectedInvoked: true,
			expectedStatus:  http.StatusOK,
			sessionKey:      "sessionKeyValue",
			val:             `{"sessionValue":{"value":"test"}}`,
			err:             nil,
		},
	}

	for _, test := range tests {

		getSessionHandler := &GetSessionHandlerGrpc{
			testExpectationsGrpc: testExpectationsGrpc{
				handlerErr: test.err,
				handlerVal: test.val,
			},
		}

		testQueries := handlers.Queries{
			GetSession: getSessionHandler,
		}

		testApp := handlers.Application{
			Queries: testQueries,
		}

		httpSvc := service.NewHttpService(testApp)

		request := httptest.NewRequest(http.MethodPost, "/api/session", strings.NewReader(""))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		httpSvc.GetSession(response, request, test.sessionKey)

		assert.True(t, response.Code == test.expectedStatus, fmt.Sprintf("Should respond with status code %d\n", test.expectedStatus))

		if test.expectedInvoked {
			assert.True(t, getSessionHandler.invoked, "'Handle' should have been invoked")
		} else {
			assert.False(t, getSessionHandler.invoked, "'Handle' should not have been invoked")
		}
	}

}

func TestDeleteHttpSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario        string
		expectedInvoked bool
		expectedStatus  int
		sessionKey      string
		err             error
	}{
		{
			scenario:        "Should respond with bad request if SessionKey is empty",
			expectedInvoked: false,
			expectedStatus:  http.StatusBadRequest,
			sessionKey:      "",
			err:             nil,
		},
		{
			scenario:        "Should respond with internal server error if handler responds with error",
			expectedInvoked: true,
			expectedStatus:  http.StatusInternalServerError,
			sessionKey:      "sessionKeyValue",
			err:             fmt.Errorf("Error from handler"),
		},
		{
			scenario:        "Should respond with accepted if handler returns without errors",
			expectedInvoked: true,
			expectedStatus:  http.StatusAccepted,
			sessionKey:      "sessionKeyValue",
			err:             nil,
		},
	}

	for _, test := range tests {

		deleteSessionHandler := &DeleteSessionHandlerGrpc{
			testExpectationsGrpc: testExpectationsGrpc{
				handlerErr: test.err,
			},
		}

		deleteCommand := handlers.Commands{
			DeleteSession: deleteSessionHandler,
		}

		testApp := handlers.Application{
			Commands: deleteCommand,
		}

		httpSvc := service.NewHttpService(testApp)

		request := httptest.NewRequest(http.MethodPost, "/api/session", strings.NewReader(""))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		httpSvc.DeleteSession(response, request, test.sessionKey)

		assert.True(t, response.Code == test.expectedStatus, fmt.Sprintf("Should respond with status code %d\n", test.expectedStatus))

		if test.expectedInvoked {
			assert.True(t, deleteSessionHandler.invoked, "'Handle' should have been invoked")
		} else {
			assert.False(t, deleteSessionHandler.invoked, "'Handle' should not have been invoked")
		}
	}

}
