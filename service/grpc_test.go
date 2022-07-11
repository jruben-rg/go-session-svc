package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jruben-rg/go-session-svc/genproto/session"
	"github.com/jruben-rg/go-session-svc/service"
	"github.com/jruben-rg/go-session-svc/sessions/handlers"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/command"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/query"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type testExpectationsGrpc struct {
	invoked    bool
	handlerVal interface{}
	handlerErr error
}

type DeleteSessionHandlerGrpc struct {
	command.DeleteSessionHandler
	testExpectationsGrpc
}

type SetSessionHandlerGrpc struct {
	command.SetSessionHandler
	testExpectationsGrpc
}

type GetSessionHandlerGrpc struct {
	query.GetSessionHandler
	testExpectationsGrpc
}

func (d *DeleteSessionHandlerGrpc) Handle(ctx context.Context, cmd command.DeleteSession) error {
	d.invoked = true
	return d.handlerErr
}

func (s *SetSessionHandlerGrpc) Handle(ctx context.Context, cmd command.SetSession) error {
	s.invoked = true
	return s.handlerErr
}

func (g *GetSessionHandlerGrpc) Handle(ctx context.Context, cmd query.GetSession) (interface{}, error) {
	g.invoked = true
	return g.handlerVal, g.handlerErr
}

func TestSetGrpcSession(t *testing.T) {
	t.Parallel()

	sessionValue, err := structpb.NewValue(map[string]interface{}{
		"test": "Grpc",
	})

	if err != nil {
		t.Errorf("Cannot create session value.")
	}

	tests := []struct {
		scenario        string
		expectedInvoked bool
		expectedStatus  codes.Code
		expectedError   bool
		sessionRequest  session.SetSessionRequest
		handlerErr      error
	}{
		{
			scenario:        "Should respond with bad request if SessionKey empty",
			expectedInvoked: false,
			expectedError:   true,
			expectedStatus:  codes.InvalidArgument,
			sessionRequest:  session.SetSessionRequest{Session: &session.Session{Key: ""}},
			handlerErr:      nil,
		},
		{
			scenario:        "Should respond with bad request if Session Value empty",
			expectedInvoked: false,
			expectedError:   true,
			expectedStatus:  codes.InvalidArgument,
			sessionRequest:  session.SetSessionRequest{Session: &session.Session{Key: "Key", Value: nil}},
			handlerErr:      nil,
		},
		{
			scenario:        "Should respond with internal error if handler returns an error",
			expectedInvoked: true,
			expectedError:   true,
			expectedStatus:  codes.Internal,
			sessionRequest:  session.SetSessionRequest{Session: &session.Session{Key: "Key", Value: sessionValue.GetStructValue()}},
			handlerErr:      fmt.Errorf("Error from handler"),
		},
		{
			scenario:        "Should not return any errors if no errors are found",
			expectedInvoked: true,
			expectedError:   false,
			sessionRequest:  session.SetSessionRequest{Session: &session.Session{Key: "Key", Value: sessionValue.GetStructValue()}},
			handlerErr:      nil,
		},
	}

	for _, test := range tests {

		setSessionHandler := &SetSessionHandlerGrpc{
			testExpectationsGrpc: testExpectationsGrpc{
				handlerErr: test.handlerErr,
			},
		}

		setCommand := handlers.Commands{
			SetSession: setSessionHandler,
		}

		appSet := handlers.Application{
			Commands: setCommand,
		}

		grpcSvc := service.NewGrpcService(appSet)

		_, err := grpcSvc.SetSession(context.Background(), &test.sessionRequest)

		if test.expectedError {

			assert.True(t, err != nil, fmt.Sprintf("Was expecting an error for scenario %s\n", test.scenario))

			if err != nil {
				e, ok := status.FromError(err)
				if !ok {
					t.Errorf("Non grpc error for scenario %s\n", test.scenario)
				}

				e.Code()
				assert.True(t, e.Code() == test.expectedStatus, fmt.Sprintf("Expected error is '%d', found '%d'. Scenario %s\n", test.expectedStatus, e.Code(), test.scenario))
			}

		} else {
			assert.Nil(t, err, fmt.Sprintf("Wasnt expecting an error for scenario '%s'. Got '%v'.\n", test.scenario, err))
		}

		if test.expectedInvoked {
			assert.True(t, setSessionHandler.invoked, "'Handle' should have been invoked")
		} else {
			assert.False(t, setSessionHandler.invoked, "'Handle' should not have been invoked")
		}
	}

}

func TestGetGrpcSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario        string
		expectedStatus  codes.Code
		expectedError   bool
		sessionRequest  session.GetSessionRequest
		handlerInvoked  bool
		handlerErr      error
		handlerResponse interface{}
	}{
		{
			scenario:       "Should respond with Invalid Argument if SessionKey empty",
			expectedError:  true,
			expectedStatus: codes.InvalidArgument,
			sessionRequest: session.GetSessionRequest{Key: ""},
			handlerInvoked: false,
			handlerErr:     nil,
		},
		{
			scenario:       "Should respond with Internal error if handler returns an error",
			expectedError:  true,
			expectedStatus: codes.Internal,
			sessionRequest: session.GetSessionRequest{Key: "Key"},
			handlerInvoked: true,
			handlerErr:     fmt.Errorf("Error from handler"),
		},
		{
			scenario:        "Should respond Not Found if handler returns empty response",
			expectedError:   true,
			expectedStatus:  codes.NotFound,
			sessionRequest:  session.GetSessionRequest{Key: "Key"},
			handlerInvoked:  true,
			handlerErr:      nil,
			handlerResponse: "",
		},
		{
			scenario:        "Should return session value",
			expectedError:   false,
			sessionRequest:  session.GetSessionRequest{Key: "Key"},
			handlerInvoked:  true,
			handlerErr:      nil,
			handlerResponse: `{"response":"value"}`,
		},
	}

	for _, test := range tests {

		getSessionHandler := &GetSessionHandlerGrpc{
			testExpectationsGrpc: testExpectationsGrpc{
				handlerErr: test.handlerErr,
				handlerVal: test.handlerResponse,
			},
		}

		getQuery := handlers.Queries{
			GetSession: getSessionHandler,
		}

		appSet := handlers.Application{
			Queries: getQuery,
		}

		grpcSvc := service.NewGrpcService(appSet)

		sessionResponse, err := grpcSvc.GetSession(context.Background(), &test.sessionRequest)

		if test.expectedError {

			assert.True(t, err != nil, fmt.Sprintf("Was expecting an error for scenario %s\n", test.scenario))

			if err != nil {
				e, ok := status.FromError(err)
				if !ok {
					t.Errorf("Non grpc error for scenario %s\n", test.scenario)
				}

				e.Code()
				assert.True(t, e.Code() == test.expectedStatus, fmt.Sprintf("Expected error is '%d', found '%d'. Scenario %s\n", test.expectedStatus, e.Code(), test.scenario))
			}

		} else {
			fmt.Printf("sessionResponse %v\n", sessionResponse.Session.Value.String())
			assert.True(t, test.sessionRequest.Key == sessionResponse.Session.Key, "Session Key should match")
			assert.Nil(t, err, fmt.Sprintf("Wasnt expecting an error for scenario '%s'. Got '%v'.\n", test.scenario, err))
		}

		if test.handlerInvoked {
			assert.True(t, getSessionHandler.invoked, "'Handle' should have been invoked")
		} else {
			assert.False(t, getSessionHandler.invoked, "'Handle' should not have been invoked")
		}
	}

}

func TestDeleteGrpcSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		scenario        string
		expectedInvoked bool
		expectedStatus  codes.Code
		expectedError   bool
		sessionRequest  session.DeleteSessionRequest
		handlerErr      error
	}{
		{
			scenario:        "Should respond with bad request if SessionKey empty",
			expectedInvoked: false,
			expectedError:   true,
			expectedStatus:  codes.InvalidArgument,
			sessionRequest:  session.DeleteSessionRequest{Key: ""},
			handlerErr:      nil,
		},
		{
			scenario:        "Should respond with bad request if handler returns an error",
			expectedInvoked: true,
			expectedError:   true,
			expectedStatus:  codes.Internal,
			sessionRequest:  session.DeleteSessionRequest{Key: "Key"},
			handlerErr:      fmt.Errorf("Error from handler"),
		},
		{
			scenario:        "Should respond without errors if handler does not return an error",
			expectedInvoked: true,
			expectedError:   false,
			sessionRequest:  session.DeleteSessionRequest{Key: "Key"},
			handlerErr:      nil,
		},
	}

	for _, test := range tests {

		deleteSessionHandler := &DeleteSessionHandlerGrpc{
			testExpectationsGrpc: testExpectationsGrpc{
				handlerErr: test.handlerErr,
			},
		}

		deleteCommand := handlers.Commands{
			DeleteSession: deleteSessionHandler,
		}

		appSet := handlers.Application{
			Commands: deleteCommand,
		}

		grpcSvc := service.NewGrpcService(appSet)

		_, err := grpcSvc.DeleteSession(context.Background(), &test.sessionRequest)

		if test.expectedError {

			assert.True(t, err != nil, fmt.Sprintf("Was expecting an error for scenario %s\n", test.scenario))

			if err != nil {
				e, ok := status.FromError(err)
				if !ok {
					t.Errorf("Non grpc error for scenario %s\n", test.scenario)
				}

				e.Code()
				assert.True(t, e.Code() == test.expectedStatus, fmt.Sprintf("Expected error is '%d', found '%d'. Scenario %s\n", test.expectedStatus, e.Code(), test.scenario))
			}

		} else {
			assert.Nil(t, err, fmt.Sprintf("Wasnt expecting an error for scenario '%s'. Got '%v'.\n", test.scenario, err))
		}

		if test.expectedInvoked {
			assert.True(t, deleteSessionHandler.invoked, "'Handle' should have been invoked")
		} else {
			assert.False(t, deleteSessionHandler.invoked, "'Handle' should not have been invoked")
		}
	}

}
