.PHONY: run

run:
	go run cmd/main.go

.PHONY: migrate

migrate:
	go run cmd/migrator/migrator.go

generate:
	oapi-codegen -config internal/clients/wowaudit/codegen.yaml internal/clients/wowaudit/codegen.yaml