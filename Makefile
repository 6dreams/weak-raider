.PHONY: run

run:
	go run cmd/main.go -config config.yaml -env .env

.PHONY: migrate

migrate:
	go run cmd/migrator/migrator.go -config config.yaml -env .env

gen:
	oapi-codegen -config internal/clients/wowaudit/codegen.yaml internal/clients/wowaudit/wowaudit.yaml
	oapi-codegen -config internal/clients/blizzard/codegen.yaml internal/clients/blizzard/blizzard.yaml
	oapi-codegen -config internal/clients/warcraftlogs/codegen.yaml internal/clients/warcraftlogs/warcraftlogs.yaml
