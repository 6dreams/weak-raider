# weak-raider
Tooling for World of Warcraft.

# Development
## Инициализация проекта
```bash
cp .env.dist .env
go install github.com/oapi-codegen/oapi-codegen
go get github.com/Khan/genqlient
```

Запустить docker-compose (если база будет в docker):
```bash
docker-compose up
```

## Запуск приложения
```bash
make run
```

## Получение изменённой схемы для применения миграций
```bash
make schema-diff
```

# Цели
* Проверка талантов рейдеров и сравнение их с часто используемыми.
* Оценка требуемых предметов, какие предметы с рейда или подземелий дадут +урон и сколько

# Рабочий процесс
## Синхронизация сезонов, подземелий
BlizzAPI -> DB
