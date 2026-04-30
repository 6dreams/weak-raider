# weak-raider
Tooling for World of Warcraft.

# Development
```bash
cp .env.dist .env
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
