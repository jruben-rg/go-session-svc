package command

import (
	"context"

	"github.com/jruben-rg/go-commons-handler/decorator"
	"github.com/jruben-rg/go-session-svc/sessions/domain/session"
	"github.com/sirupsen/logrus"
)

type DeleteSession struct {
	Key string
}

type DeleteSessionHandler decorator.CommandHandler[DeleteSession]

type deleteSessionHandler struct {
	sessionRepo session.Repository
}

func NewDeleteSessionHandler(
	sessionRepo session.Repository,
	logger *logrus.Entry,
) DeleteSessionHandler {

	if sessionRepo == nil {
		panic("nil sessionRepo")
	}

	return decorator.WithCommandDecorator[DeleteSession](
		deleteSessionHandler{sessionRepo: sessionRepo},
		logger,
	)
}

func (h deleteSessionHandler) Handle(ctx context.Context, cmd DeleteSession) error {
	_, err := h.sessionRepo.Delete(ctx, cmd.Key)
	return err
}
