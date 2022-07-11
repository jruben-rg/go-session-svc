package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jruben-rg/go-session-svc/app"
	"github.com/jruben-rg/go-session-svc/genproto/session"
	"github.com/jruben-rg/go-session-svc/server"
	"github.com/jruben-rg/go-session-svc/service"
	"google.golang.org/grpc"
)

func main() {

	application := app.NewApplication("redis")
	serverType := strings.ToLower(os.Getenv("SERVER_TYPE"))
	switch serverType {
	case "http":
		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return server.HandlerFromMux(
				service.NewHttpService(application),
				router,
			)
		})

	case "grpc":
		server.RunGrpcServer(func(server *grpc.Server) {
			svc := service.NewGrpcService(application)
			session.RegisterSessionServiceServer(server, svc)
		})

	default:
		panic(fmt.Sprintf("server type '%s' not supported", serverType))
	}

}
