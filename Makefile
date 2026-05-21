.PHONY: run

run:
	go run cmd/main.go -config=config.yml -env=.env

.PHONY: migrate-up
.PHONY: migrate-down
.PHONY: schema_diff
.PHONY: gen


migrate-up:
	go run cmd/migrator/migrator.go -config config.yml -env .env

migrate-down:
	go run cmd/migrator/migrator.go -config config.yml -env .env -migration=false

schema-diff:
	go run cmd/schema_diff/cli.go -config config.yml -env .env

gen:
	oapi-codegen -config internal/clients/wowaudit/codegen.yaml internal/clients/wowaudit/wowaudit.yaml
	oapi-codegen -config internal/clients/blizzard/codegen.yaml internal/clients/blizzard/blizzard.yaml
	oapi-codegen -config internal/clients/warcraftlogs/codegen.yaml internal/clients/warcraftlogs/warcraftlogs.yaml
	oapi-codegen -config internal/clients/raidbots/codegen.yaml internal/clients/raidbots/client.yaml
