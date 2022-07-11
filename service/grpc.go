package service

import (
	"context"
	"encoding/json"

	"github.com/jruben-rg/go-session-svc/genproto/session"
	"github.com/jruben-rg/go-session-svc/sessions/handlers"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/command"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type GrpcService struct {
	app handlers.Application
}

func NewGrpcService(application handlers.Application) GrpcService {
	return GrpcService{application}
}

func (g GrpcService) SetSession(ctx context.Context, request *session.SetSessionRequest) (*emptypb.Empty, error) {

	if request.Session.Key == "" || request.Session.Value == nil {
		return nil, status.Error(codes.InvalidArgument, "SessionKey cannot be empty")
	}

	if err := g.app.Commands.SetSession.Handle(ctx,
		command.SetSession{
			Key:   request.Session.Key,
			Value: request.Session.Value.AsMap(),
		}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcService) GetSession(ctx context.Context, request *session.GetSessionRequest) (*session.GetSessionResponse, error) {

	if request.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "SessionKey cannot be empty")
	}

	res, err := g.app.Queries.GetSession.Handle(ctx, query.GetSession{Key: request.Key})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if res == "" {
		return nil, status.Errorf(codes.NotFound, "session id not found")
	}

	resStr, ok := res.(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "cannot parse session value")
	}

	sessionValue := command.SessionValue{}
	if err := json.Unmarshal([]byte(resStr), &sessionValue); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot unmarshall session value to JSON")
	}

	structSession, err := structpb.NewStruct(sessionValue)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot transform session value to proto struct type")
	}

	return &session.GetSessionResponse{
		Session: &session.Session{
			Key:   request.Key,
			Value: structSession,
		},
	}, nil
}

func (g GrpcService) DeleteSession(ctx context.Context, request *session.DeleteSessionRequest) (*emptypb.Empty, error) {

	if request.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "SessionKey cannot be empty")
	}

	if err := g.app.Commands.DeleteSession.Handle(ctx, command.DeleteSession{Key: request.Key}); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
