package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jruben-rg/go-commons-handler/decorator"
	"github.com/jruben-rg/go-session-svc/sessions/domain/session"
	"github.com/sirupsen/logrus"
)

type SessionValue map[string]interface{}

type SetSession struct {
	Key   string
	Value SessionValue
}

type SetSessionHandler decorator.CommandHandler[SetSession]

type setSessionHandler struct {
	sessionRepo session.Repository
}

func NewSetSessionHandler(
	sessionRepo session.Repository,
	logger *logrus.Entry,
) SetSessionHandler {

	if sessionRepo == nil {
		panic("nil sessionRepo")
	}

	return decorator.WithCommandDecorator[SetSession](
		setSessionHandler{sessionRepo: sessionRepo},
		logger,
	)
}

func (h setSessionHandler) Handle(ctx context.Context, cmd SetSession) error {

	err := h.sessionRepo.Set(ctx, cmd.Key, cmd.Value)
	if err != nil {
		return fmt.Errorf("error when trying to set session %s", cmd.Key)
	}

	return nil
}

func (s SessionValue) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s SessionValue) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}
