# DockerControllerGolang

## Описание

Проект для копирования докер-контейнеров из репозиториев и их контроля

## Структура проекта

- `repo/` — репозитории
- `models/` — модели
- `db/` — модуль работы с бд
- `core/config/` — конфиг
- `api/` — роутеры
- `adapters/` — адаптеры

## Установка

```bash
git clone https://github.com/Dedushka-Lenin/DockerControllerGolang
cd DockerControllerGolang
```

## Настройка

1. swag init -g ./cmd/main.go

2. Укажите настройки в файле 'app/core/config/config.json'

## Работа программы

1. Запуск докера — ```bash sudo systemctl start docker```

2. Запуск api — ```bash go run ./cmd``

3. `http://localhost:8080/`

4. Тестовые ссылки:
    `https://github.com/Dedushka-Lenin/Hello-World-Container`