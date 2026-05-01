# weak-raider
Tooling for World of Warcraft.

# Development
```bash
cp .env.dist .env
go install github.com/oapi-codegen/oapi-codegen
make gen
```

Run docker-compose:
```bash
docker-compose up
```

Run migrations:
```bash
go run ./cmd/migrator/migrator.go
```

Run app:
```bash
go run ./cmd/main.go
```
# Цели
* Проверка талантов рейдеров и сравнение их с часто используемыми.
* Оценка требуемых предметов, какие предметы с рейда или подземелий дадут +урон и сколько

# Рабочий процесс
## Синхронизация сезонов, подземелий
BlizzAPI -> DB
