# go-session-svc
This microservice is intended to be used as a session storage service. 
The session is stored in JSON format.

The following is an example of the data that could be set using the `POST - Http` method:

```
{
    "sessionKey": "someSessionKey",
    "sessionValue": {
        "data": "someData",
        "age": 35,
        "props": {
            "someProp": "someValue"
        }
    }
}
```

where `sessionKey` is a string value, and `sessionValue` is of type `map[string]interface{}`

Default underlying memory db is `Redis`.

# Http

Provides a simple api with three methods. For additional information see file at api/openapi/session.yml

- `POST /api/session`: Stores a JSON value in memory.
- `GET /api/session/{sessionId}`: Retrieves a previously stored value
- `DELETE /api/session/{sessionId}`: Deletes an stored value

# Grpc

It offers the following RPC methods:

- `SetSession` 
- `GetSession`
- `DeleteSession`

For more info see file at api/protobuf/session.proto

# Required Config
Configuration is passed to the app by using the following environment variables:

- `SERVER_TYPE`: Must be `http` | `grpc`
- `SERVER_PORT`: Port in which the app listens
- `MEMORY_DB_HOST`: DB Host
- `MEMORY_DB_PORT`: DB Port
- `MEMORY_DB_ID`: DB Instance (Currently used for Redis DB ID)
- `MEMORY_DB_PASSWORD`: DB Password

# Docker

When using the `docker-compose.yml` file provided, it will start the following containers:
- `redis`: Underlying memory db storage
- `redis-commander`: Web interface to visualize stored sessions
- `session-http`: HTTP functionality provided by this app
- `session-grpc`: GRPC functionality provided by this app


# How to use it:
This app includes a `Makefile` to simplify all the commands whenever possible.

- `make servers`: Generates `openapi` and `grpc` server, client and required types.
  To generate required files, install [oapi-codegen](https://github.com/deepmap/oapi-codegen) and [protoc](https://grpc.io/docs/languages/go/quickstart/)  
- `make docker-up`: Start Docker services specified in the section above.
- `make docker-down`: Stop Docker services