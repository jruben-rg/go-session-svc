package app

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jruben-rg/go-session-svc/sessions/adapters"
	"github.com/jruben-rg/go-session-svc/sessions/domain/session"
	"github.com/jruben-rg/go-session-svc/sessions/handlers"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/command"
	"github.com/jruben-rg/go-session-svc/sessions/handlers/query"
	"github.com/sirupsen/logrus"
)

const (
	aDay  = 60 * 60 * 24
	aYear = aDay * 365
)

func NewApplication(dbType string) handlers.Application {

	var sessionRepo session.Repository
	logger := logrus.NewEntry(logrus.StandardLogger())
	var redisDb int

	durStr := getEnvVar("MEMORY_DB_DURATION", fmt.Sprint(aDay))
	dur := toInt(durStr)

	duration := time.Duration(dur) * time.Second

	host := getEnvVar("MEMORY_DB_HOST", "localhost")
	port := getEnvVar("MEMORY_DB_PORT", "6379")
	password := getEnvVar("MEMORY_DB_PASSWORD", "")
	addr := fmt.Sprintf("%s:%s", host, port)

	switch dbType {
	case "redis":
		db := getEnvVar("MEMORY_DB_ID", "0")
		redisDb = toInt(db)
		sessionRepo = adapters.NewRedisCache(addr, redisDb, password, duration)
	default:
		panic(fmt.Sprintf("db type '%s' not supported", dbType))
	}

	return handlers.Application{
		Commands: handlers.Commands{
			DeleteSession: command.NewDeleteSessionHandler(sessionRepo, logger),
			SetSession:    command.NewSetSessionHandler(sessionRepo, logger),
		},
		Queries: handlers.Queries{
			GetSession: query.NewGetSessionHandler(sessionRepo, logger),
		},
	}
}

func getEnvVar(varName, varDefaultValue string) string {
	varValue := os.Getenv(varName)
	if varValue == "" {
		return varDefaultValue
	}

	return varValue
}

func toInt(valueStr string) int {
	dur, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Errorf("error '%s' when parsing to integer value", err))
	}
	return dur
}
