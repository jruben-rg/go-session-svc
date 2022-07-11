FROM golang:1.18-alpine AS build

WORKDIR /src/go-session-svc

ADD app /src/go-session-svc/app
ADD genproto /src/go-session-svc/genproto
ADD server /src/go-session-svc/server
ADD service /src/go-session-svc/service
ADD sessions /src/go-session-svc/sessions
COPY go.* /src/go-session-svc/
COPY main.go /src/go-session-svc/
RUN CGO_ENABLED=0 go build -o /bin/go-session-svc

ENTRYPOINT [ "/bin/go-session-svc" ]