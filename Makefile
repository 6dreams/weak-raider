.PHONY: run

run:
	go run cmd/main.go -config=config.yml -env=.env

.PHONY: migrate
.PHONY: gen


migrate:
	go run cmd/migrator/migrator.go -config config.yml -env .env

schema-diff:
	go run cmd/schema_diff/cli.go -config config.yml -env .env

gen:
	oapi-codegen -config internal/clients/wowaudit/codegen.yaml internal/clients/wowaudit/wowaudit.yaml
	oapi-codegen -config internal/clients/blizzard/codegen.yaml internal/clients/blizzard/blizzard.yaml
	oapi-codegen -config internal/clients/warcraftlogs/codegen.yaml internal/clients/warcraftlogs/warcraftlogs.yaml
