# 🎵 Music Platform API

Backend-сервис музыкальной платформы.

---

### ✨ Возможности

- 🎧 **Каталог треков (публичный)**
    - Список треков с пагинацией
    - Поиск по названию, исполнителю, альбому и жанру
    - Справочник жанров
    - Прослушивание трека с записью в историю
- 📻 **Плейлисты (только свои)**
    - Создание, просмотр, обновление, удаление
    - Добавление и удаление треков
    - Лимит плейлистов для FREE-подписки
- ❤️ **Избранное и история**
    - Добавление и удаление треков в избранном (лимит для FREE)
    - История прослушиваний пользователя
- 💎 **Подписки**
    - Тарифы FREE / PREMIUM
    - Конфигурируемые лимиты FREE (`FREE_PLAYLIST_LIMIT`, `FREE_FAVORITES_LIMIT`)
- 🔐 **Аккаунты и доступ**
    - Регистрация, вход и выход
    - Роли USER / ADMIN
- 🛡️ **Админ-панель (роль ADMIN)**
    - Загрузка треков в каталог
    - Обновление и удаление треков
    - Управление подпиской пользователей

---

### 🛠 Технологический стек

| Категория | Технологии |
|---|---|
| **Язык** | Go 1.26+ |
| **Транспорт** | REST (chi v5), OpenAPI 3.1, Swagger UI |
| **БД** | PostgreSQL 18, pgx v5 (пул), sqlc (генерация запросов) |
| **Миграции** | goose (автоприменение при старте) |
| **Кэш** | Redis 8.0 (go-redis v9) |
| **Хранилище** | RustFS (S3-совместимое), minio-go v7 |
| **Аутентификация** | JWT (golang-jwt v5, Access + Refresh), bcrypt |
| **Валидация** | go-playground/validator v10 |
| **Конфигурация** | caarlos0/env v11 + godotenv |
| **Логирование** | Zap (структурированный JSON) |
| **Тестирование** | testify, httptest, интеграционные тесты репозиториев |
| **Контейнеризация** | Docker, Docker Compose |

---

### 🏗 Архитектура

```
.
├── api/v1/                                  # Спецификация OpenAPI 3.1
│   ├── openapi.yaml                         #   Корневой файл спецификации
│   ├── paths/                               #   Эндпоинты по доменам
│   │   ├── auth.yaml
│   │   ├── users.yaml
│   │   ├── tracks.yaml
│   │   ├── playlists.yaml
│   │   ├── favorites.yaml
│   │   ├── history.yaml
│   │   └── admin.yaml
│   └── components/                          #   Схемы, параметры, ответы, security
│       ├── schemas/
│       ├── parameters.yaml
│       ├── responses.yaml
│       └── securitySchemes.yaml
├── cmd/
│   └── api/main.go                          # Точка входа приложения
├── database/
│   ├── migrations/                          # SQL-миграции goose (+ seed-данные)
│   │   ├── 00001_create_users.sql
│   │   ├── 00002_enable_pg_trgm.sql
│   │   ├── ...
│   │   └── 00011_seed_data.sql
│   └── queries/                             # SQL-запросы для sqlc
│       ├── catalog.sql
│       ├── favorites.sql
│       ├── history.sql
│       ├── playlists.sql
│       ├── tracks.sql
│       └── users.sql
├── internal/
│   ├── adapter/
│   │   ├── cache/redis/                     # Клиент Redis
│   │   ├── database/postgres/               # Пул pgx, миграции, sqlc-код
│   │   │   └── sqlc/                        #   Сгенерированные запросы
│   │   ├── storage/s3/                      # Клиент S3 (minio-go)
│   │   └── transport/http/
│   │       ├── middleware/                  # CORS, auth, role, logger, recover, request_id
│   │       ├── swagger/                     # Раздача Swagger UI
│   │       ├── v1/                          # Хендлеры по доменам
│   │       │   ├── auth/
│   │       │   ├── user/
│   │       │   ├── track/
│   │       │   ├── playlist/
│   │       │   ├── favorite/
│   │       │   ├── history/
│   │       │   └── admin/
│   │       └── router.go                    # Сборка роутера chi
│   ├── app/                                 # Bootstrap и DI-контейнер
│   │   ├── app.go                           #   Запуск + graceful shutdown
│   │   └── container.go                     #   Конфиг → клиенты → репо → кэш → сервисы
│   ├── cache/                               # Слой кэширования (Redis)
│   │   ├── blacklist/                       #   Отозванные access-токены
│   │   ├── refresh/                         #   Allowlist refresh-токенов
│   │   ├── track/ genre/ search/ popular/   #   Кэш каталога
│   │   └── caches.go
│   ├── config/                              # Конфигурация (env + валидация)
│   ├── domain/                              # Доменные модели и ошибки
│   │   ├── user/                            #   User, Role, Subscription, Password
│   │   ├── track/
│   │   ├── playlist/
│   │   ├── favorite/
│   │   └── history/
│   ├── repository/                          # Слой доступа к данным (PostgreSQL)
│   │   ├── user/ track/ playlist/
│   │   ├── favorite/ history/
│   │   └── testutil/                        #   Хелперы интеграционных тестов
│   └── service/                             # Бизнес-логика
│       ├── auth/ user/ track/ playlist/
│       ├── favorite/ history/ admin/
│       └── services.go
├── pkg/                                     # Переиспользуемые пакеты
│   ├── httpx/                               #   JSON/error-хелперы, context, Bearer
│   ├── jwt/                                 #   JWT-менеджер (claims, подпись, парсинг)
│   └── logger/                              #   Фабрика zap-логгера
├── scripts/
│   └── seed-s3.sh                           # Сидинг S3-бакета аудиофайлами
├── web/swagger/                             # Swagger UI (index.html)
├── docker-compose.yml                       # API + PostgreSQL + Redis + RustFS
├── Dockerfile                               # Multi-stage сборка
├── sqlc.yaml
└── .env.example
```

#### Ключевые принципы

- **Слоистая архитектура (Clean Architecture).** Запрос проходит путь `Handler → Service → Repository`. Каждый слой знает только о соседнем снизу: хендлеры не трогают БД, репозитории не содержат бизнес-логики. Зависимости направлены внутрь — к домену.
- **Домен в центре.** Бизнес-модели, правила (роли, подписки, пароли) и доменные ошибки живут в `internal/domain` и не зависят ни от транспорта, ни от БД. Хендлеры маппят доменные ошибки в HTTP-коды на своей стороне.
- **Интерфейсы на стороне потребителя.** Не репозиторий объявляет, что он умеет, а сервис описывает, какой репозиторий ему нужен (`TrackRepository`, `TrackCache` в `internal/service/track`). Хендлер так же описывает нужный ему `Service`. Это развязывает слои и упрощает мок-тестирование.
- **Тонкие хендлеры.** Схема каждого хендлера: parse → validate → вызов сервиса → запись ответа. Вся логика (включая проверку лимитов FREE) — в сервисном слое, потому что ей нужны данные из БД.
- **Dependency Injection.** Весь граф зависимостей собирается в одном месте — `internal/app/container.go`: конфиг → клиенты (PostgreSQL, Redis, S3) → репозитории → кэши → сервисы → роутер. Никаких глобальных переменных и `init()`.
- **`context.Context` сквозь все слои.** Контекст пробрасывается от middleware до запроса в БД — таймауты и отмена запроса работают на всей глубине стека.
- **Cache-aside.** Кэш — отдельный слой (`internal/cache`), сервис сначала идёт в Redis, при промахе — в БД с последующей записью в кэш; при изменении данных админом кэш инвалидируется.
- **Безопасность по умолчанию.** Пароли — bcrypt + соль, секреты — только из env, токены и пароли не пишутся в логи, soft delete вместо физического удаления треков.

---

### 📡 API эндпоинты

#### Аутентификация (`/api/v1/auth`)

| Метод | Эндпоинт | Описание |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Регистрация (роль USER, подписка FREE) |
| `POST` | `/api/v1/auth/login` | Вход, выдача access + refresh токенов |
| `POST` | `/api/v1/auth/refresh` | Ротация access & refresh токенов |
| `POST` | `/api/v1/auth/logout` | Выход (blacklist access, отзыв refresh) |

#### Профиль (`/api/v1/users`)

| Метод | Эндпоинт | Описание |
|---|---|---|
| `GET` | `/api/v1/users/me` | Профиль текущего пользователя |
| `PATCH` | `/api/v1/users/me` | Обновление профиля |

#### Треки (`/api/v1/tracks`)

| Метод | Эндпоинт | Описание |
|---|---|---|
| `GET` | `/api/v1/tracks` | Список треков (пагинация `limit`/`offset`) |
| `GET` | `/api/v1/tracks/search?q=...` | Поиск по названию/исполнителю/альбому/жанру |
| `GET` | `/api/v1/tracks/genres` | Список жанров |
| `GET` | `/api/v1/tracks/{id}` | Трек по ID |
| `POST` | `/api/v1/tracks/{id}/play` | Прослушать трек (запись в историю) |

#### Плейлисты (`/api/v1/playlists`)

| Метод | Эндпоинт | Описание |
|---|---|---|
| `POST` | `/api/v1/playlists` | Создать плейлист _(FREE: лимит)_ |
| `GET` | `/api/v1/playlists` | Мои плейлисты |
| `GET` | `/api/v1/playlists/{id}` | Получить плейлист |
| `PUT` | `/api/v1/playlists/{id}` | Обновить плейлист |
| `DELETE` | `/api/v1/playlists/{id}` | Удалить плейлист |
| `GET` | `/api/v1/playlists/{id}/tracks` | Треки в плейлисте |
| `POST` | `/api/v1/playlists/{playlist_id}/tracks/{track_id}` | Добавить трек |
| `DELETE` | `/api/v1/playlists/{playlist_id}/tracks/{track_id}` | Удалить трек |

#### Избранное (`/api/v1/favorites`)

| Метод | Эндпоинт | Описание |
|---|---|---|
| `GET` | `/api/v1/favorites/tracks` | Список избранного |
| `POST` | `/api/v1/favorites/tracks/{track_id}` | Добавить в избранное _(FREE: лимит)_ |
| `DELETE` | `/api/v1/favorites/tracks/{track_id}` | Удалить из избранного |

#### История прослушиваний (`/api/v1/listening-history`)

| Метод | Эндпоинт | Описание |
|---|---|---|
| `GET` | `/api/v1/listening-history` | История прослушиваний пользователя |

#### Админ (`/api/v1/admin`) — только роль ADMIN

| Метод | Эндпоинт | Описание |
|---|---|---|
| `POST` | `/api/v1/admin/tracks` | Загрузить трек (`multipart/form-data` → S3 + БД) |
| `PUT` | `/api/v1/admin/tracks/{id}` | Обновить трек (с инвалидацией кэша) |
| `DELETE` | `/api/v1/admin/tracks/{id}` | Удалить трек (soft delete) |
| `PATCH` | `/api/v1/admin/users/{id}/subscription` | Изменить подписку пользователя |

#### Служебные

| Метод | Эндпоинт | Описание |
|---|---|---|
| `GET` | `/health` | Проверка живости сервиса |

> **Политика доступа:** регистрация, вход, refresh и весь каталог треков (список, поиск, жанры, трек по ID) — **публичные**.
> Профиль, плейлисты, избранное, история, logout и `POST /tracks/{id}/play` требуют **аутентификации** (любая роль).
> Эндпоинты `/admin/*` доступны только роли **ADMIN**.

---

### 🌐 Веб-интерфейсы

| Сервис | URL | Описание |
|---|---|---|
| **Swagger UI** | `http://localhost:8080/swagger` | Интерактивная документация и тестирование API |
| **RustFS Console** | `http://localhost:9001` | Веб-консоль S3-хранилища (бакет `tracks`) |

---

### 🏗 Установка и запуск

#### Требования

- [Docker и Docker Compose](https://docs.docker.com/engine/install/)
- [Go 1.26+](https://go.dev/dl/) — только для локального запуска без Docker

#### 1. Клонировать репозиторий

```bash
git clone https://github.com/DaniilKalts/music-platform-api.git
cd music-platform-api
```

#### 2. Настроить переменные окружения

```bash
cp .env.example .env
```

Полный справочник переменных со значениями по умолчанию:

```bash
# ─── Приложение ───────────────────────────────────────────────
APP_PORT=8080                        # Порт HTTP-сервера
SERVER_HTTP_TIMEOUT=15s              # Read/Write/Idle таймаут сервера
SERVER_HANDLER_TIMEOUT=10s           # Таймаут обработки запроса
SERVER_SHUTDOWN_TIMEOUT=15s          # Окно graceful shutdown
CORS_ALLOWED_ORIGINS=http://localhost:3000   # Разрешённые origin'ы (через запятую)

# ─── PostgreSQL ───────────────────────────────────────────────
DB_HOST=localhost                    # Хост БД ("postgres" в Docker)
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=music_platform
DB_SSL_MODE=disable
DB_MIN_CONNS=1                       # Минимум соединений в пуле
DB_MAX_CONNS=10                      # Максимум соединений в пуле
DB_MAX_CONN_LIFETIME=1h
DB_MAX_CONN_IDLE_TIME=30m
DB_STATEMENT_TIMEOUT=3s              # Таймаут одного запроса

# ─── JWT ──────────────────────────────────────────────────────
JWT_ACCESS_SECRET=access-secret-change-me    # Секрет access-токенов
JWT_REFRESH_SECRET=refresh-secret-change-me  # Секрет refresh-токенов
JWT_ACCESS_TTL=15m                   # TTL access-токена
JWT_REFRESH_TTL=720h                 # TTL refresh-токена (30 дней)

# ─── Redis ────────────────────────────────────────────────────
REDIS_HOST=localhost                 # Хост Redis ("redis" в Docker)
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s

# ─── S3 (RustFS) ─────────────────────────────────────────────
S3_ENDPOINT=localhost:9000           # Эндпоинт S3 ("rustfs:9000" в Docker)
S3_ACCESS_KEY=admin
S3_SECRET_KEY=password
S3_BUCKET=tracks                     # Бакет с аудиофайлами
S3_USE_SSL=false

# ─── Логирование ─────────────────────────────────────────────
LOG_LEVEL=info                       # debug, info, warn, error
LOG_FORMAT=json                      # json, console

# ─── Лимиты FREE-подписки ────────────────────────────────────
FREE_PLAYLIST_LIMIT=3                # Максимум плейлистов
FREE_FAVORITES_LIMIT=20              # Максимум треков в избранном
```

#### 3. Запустить через Docker Compose

```bash
docker-compose up -d --build
```

Поднимутся **4 сервиса**: API, PostgreSQL 18, Redis 8.0 и RustFS (S3).
Миграции БД (включая seed-данные: жанры, исполнители, альбомы, треки) применяются автоматически при старте API.

#### 4. Наполнить S3 аудиофайлами

```bash
./scripts/seed-s3.sh
```

Скрипт скачает public-domain записи (Vivaldi, Beethoven, Mozart, Bach, Pachelbel) с Internet Archive и загрузит их в бакет `tracks` — ключи файлов совпадают с seed-данными миграций.

#### 5. Проверить, что всё работает

| Сервис | URL | Описание |
|---|---|---|
| **REST API** | `http://localhost:8080` | HTTP API (JSON) |
| **Swagger UI** | `http://localhost:8080/swagger` | Интерактивная документация |
| **Health check** | `http://localhost:8080/health` | Живость сервиса |
| **RustFS Console** | `http://localhost:9001` | Консоль S3 (admin/password) |

```bash
# Health check
curl http://localhost:8080/health

# Регистрация пользователя
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "email": "john@example.com", "password": "Secret123"}'

# Каталог треков (публичный)
curl http://localhost:8080/api/v1/tracks

# Поиск
curl "http://localhost:8080/api/v1/tracks/search?q=vivaldi"
```

---

### 🏃 Локальный запуск (без Docker)

#### 1. Поднять инфраструктуру

Нужны запущенные PostgreSQL (`localhost:5432`), Redis (`localhost:6379`) и RustFS/MinIO (`localhost:9000`). Проще всего поднять только их через Compose:

```bash
docker-compose up -d postgres redis rustfs
```

#### 2. Настроить `.env`

В `.env.example` хосты уже указывают на `localhost` — достаточно скопировать:

```bash
cp .env.example .env
```

#### 3. Запустить приложение

```bash
go run cmd/api/main.go -config-path=.env
```

Миграции применятся автоматически при старте. Затем наполните S3:

```bash
./scripts/seed-s3.sh
```

#### 4. Тесты

```bash
# Unit-тесты (сервисы, хендлеры, middleware, jwt)
go test ./...

# Интеграционные тесты репозиториев требуют локальный PostgreSQL
# с базой music_platform_test (или переменную TEST_DATABASE_URL)
psql -U postgres -c "CREATE DATABASE music_platform_test;"
go test ./internal/repository/...
```
