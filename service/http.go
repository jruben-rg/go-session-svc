package service

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jruben-rg/go-session-svc/server"
	"github.com/jruben-rg/go-session-svc/sessions/handlers"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/command"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/query"
)

type HttpService struct {
	app handlers.Application
}

func NewHttpService(application handlers.Application) HttpService {
	return HttpService{
		app: application,
	}
}

func (h HttpService) SetSession(w http.ResponseWriter, r *http.Request) {

	postSession := server.PostSession{}
	if err := render.Decode(r, &postSession); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if postSession.SessionKey == "" {
		http.Error(w, "SessionId cannot be empty", http.StatusBadRequest)
		return
	}

	err := h.app.Commands.SetSession.Handle(r.Context(), command.SetSession{
		Key:   postSession.SessionKey,
		Value: postSession.SessionValue,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h HttpService) DeleteSession(w http.ResponseWriter, r *http.Request, sessionId string) {

	if sessionId == "" {
		http.Error(w, "SessionId cannot be emmpty", http.StatusBadRequest)
		return
	}

	err := h.app.Commands.DeleteSession.Handle(r.Context(), command.DeleteSession{
		Key: sessionId,
	})

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)

}

func (h HttpService) GetSession(w http.ResponseWriter, r *http.Request, sessionId string) {

	if sessionId == "" {
		http.Error(w, "SessionId cannot be emmpty", http.StatusBadRequest)
		return
	}

	res, err := h.app.Queries.GetSession.Handle(r.Context(), query.GetSession{
		Key: sessionId,
	})

	sessionStr, ok := res.(string)
	if !ok {
		http.Error(w, "Cannot parse session value", http.StatusInternalServerError)
		return
	}

	if sessionStr == "" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session := command.SessionValue{}
	if err := json.Unmarshal([]byte(sessionStr), &session); err != nil {
		http.Error(w, "Error when unmarshalling session value to json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	render.Respond(w, r, session)
}
