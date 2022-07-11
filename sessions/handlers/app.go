package handlers

import (
	"github.com/jruben-rg/go-session-svc/sessions/handlers/command"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/query"
)

type Commands struct {
	DeleteSession command.DeleteSessionHandler
	SetSession    command.SetSessionHandler
}

type Queries struct {
	GetSession query.GetSessionHandler
}

type Application struct {
	Commands Commands
	Queries  Queries
}
