<img width="720" alt="contestr" src="https://github.com/user-attachments/assets/58094b6c-4166-4dc4-8fbe-56a028b26e60" />

# Contestr

Веб-сервис для проведения регатных контестов: зрители следят за таблицей и ходом соревнования.

## Требования

**Обязательно:**

- [Go 1.24+](contestr/go.mod) — бэкенд
- [Node.js 18+](contestr-front/package.json) и **npm** — фронтенд
- **MongoDB**

**Для сборки бэкенда** (codegen):

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1
go install github.com/google/wire/cmd/wire@v0.6.0
```

**Опционально:** Docker и Docker Compose — запуск всего стека одной командой.

## Быстрый старт (Docker Compose)

```bash
cp contestr/.env.example contestr/.env   # заполнить секреты
docker compose up --build
```

- UI: **http://localhost:10888**
- MongoDB с хоста: `localhost:15432`

Backend с хоста: скопируйте [`docker-compose.override.example.yml`](docker-compose.override.example.yml) в `docker-compose.override.yml` — пробросится `127.0.0.1:8080:8080`.

## Локальная разработка

### MongoDB

```bash
docker compose up mongodb -d
```

Порт `15432` на хосте → `27017` в контейнере. Для локального бэкенда укажите в `contestr/.env`:

```
APP_MONGODB_URI=mongodb://localhost:15432
```

### Backend

```bash
cd contestr
cp .env.example .env
make build          # codegen + go build
./bin/app           # слушает :8080
```

Проверка: `curl http://127.0.0.1:8080/` → `HTTPServer is running`

### Frontend

```bash
cd contestr-front
npm install
npm run dev         # http://localhost:5173, прокси /api → :8080
```

Другой хост бэкенда:

```bash
VITE_API_PROXY_TARGET=http://host:8080 npm run dev
```

Подробнее о разработке фронта: [contestr-front/README.md](contestr-front/README.md), [contestr-front/DEVELOP.md](contestr-front/DEVELOP.md).

## Конфигурация

| Источник | Файл | Примеры |
|----------|------|---------|
| YAML | [`contestr/configs/config.yaml`](contestr/configs/config.yaml) | порты HTTP, интервалы sync, URL Ejudge, object storage |
| Секреты (env) | [`contestr/.env.example`](contestr/.env.example) | Codeforces API, Yandex Object Storage |

Переменные окружения с префиксом `APP_` переопределяют YAML. Путь к конфигу: `CONFIG_PATH` (по умолчанию `configs/config.yaml`).

**Для синхронизации с Codeforces** нужны `APP_CODEFORCES_API_KEY` и `APP_CODEFORCES_API_SECRET`.

**Опционально:** Object Storage — без ключей Yandex Cloud в `.env`.

## Полезные команды

| Команда | Где | Действие |
|---------|-----|----------|
| `make build` / `make test` | `contestr/` | Сборка и тесты |
| `make code-gen` | `contestr/` | OpenAPI + Wire (перед build) |
| `npm run build` / `npm run lint` | `contestr-front/` | Production-сборка и линтер |
| `npm run generate-client` | `contestr-front/` | Перегенерация API-клиента из OpenAPI |
