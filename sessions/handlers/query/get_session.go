package query

import (
	"context"

	"github.com/jruben-rg/go-commons-handler/decorator"
	"github.com/jruben-rg/go-session-svc/sessions/domain/session"
	"github.com/sirupsen/logrus"
)

type GetSession struct {
	Key string
}

type GetSessionHandler decorator.QueryHandler[GetSession, interface{}]

type getSessionHandler struct {
	sessionRepo session.Repository
}

func NewGetSessionHandler(
	sessionRepo session.Repository,
	logger *logrus.Entry,
) GetSessionHandler {

	if sessionRepo == nil {
		panic("nil SessionRepo")
	}

	return decorator.WithQueryDecorators[GetSession, interface{}](
		getSessionHandler{sessionRepo: sessionRepo},
		logger,
	)
}

func (h getSessionHandler) Handle(ctx context.Context, getSession GetSession) (interface{}, error) {
	return h.sessionRepo.Get(ctx, getSession.Key)
}
