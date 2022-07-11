include .env
export

.PHONY: servers
servers: openapi proto

.PHONY: openapi
openapi:
	@./scripts/openapi-http.sh

.PHONY: proto
proto:
	@./scripts/proto.sh

.PHONY: docker-up
docker-up:
	docker compose up -d

.PHONY: docker-down
docker-down:
	docker compose down