.PHONY: run

run:
	go run cmd/main.go

.PHONY: migrate

migrate:
	go run cmd/migrator/migrator.go

gen:
	oapi-codegen -config internal/clients/wowaudit/codegen.yaml internal/clients/wowaudit/wowaudit.yaml
	oapi-codegen -config internal/clients/blizzard/codegen.yaml internal/clients/blizzard/blizzard.yaml
