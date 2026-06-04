# TODO

Порядок сборки — сверху вниз. Отмечай пункты по мере готовности.

**Архитектура — как в `rbk-school/7-week/api-service`** (Clean Architecture, один сервис). Конвенции, нейминг пакетов и тулчейн повторяем оттуда. **Ничего сверх ТЗ не добавляем** (никакого gateway, лишних доменов, Kafka и т.п.).

> **Soft delete:** используем `deleted_at` (nullable timestamptz) в `tracks`. `is_deleted` из ТЗ игнорируем.
>
> **Аудиофайлы:** лежат в бакете RustFS (S3/MinIO). Бэкенд работает с S3 через `minio-go` SDK. Админ загружает файл через `POST /admin/tracks` (multipart/form-data).
>
> **Лимиты FREE:** прослушивания и скачивание — безлимит. Ограничены только плейлисты (`FREE_PLAYLIST_LIMIT`) и избранное (`FREE_FAVORITES_LIMIT`). Проверка лимита — в **сервисном слое** (нужен `COUNT` из БД), не в middleware.
>
> **Сквозные правила:** `context.Context` пробрасываем во все слои до БД. Интерфейсы объявляем на стороне **потребителя** (сервис описывает нужный ему репозиторий). Хендлеры тонкие: parse → validate → call service → write. БД только в репозитории. Пароли — bcrypt+salt, секреты — из env, токены/пароли не логируем.

---

## Стек и библиотеки (как в референсе)

| Назначение | Библиотека |
|---|---|
| Роутер | `go-chi/chi/v5` |
| Драйвер БД | `jackc/pgx/v5` (пул) |
| Генерация запросов | **sqlc** (`database/queries` → `internal/adapter/database/postgres/sqlc`) |
| Миграции | **goose** (`pressly/goose/v3`, файлы `00001_*.sql`) |
| Конфиг | `caarlos0/env/v11` + `joho/godotenv` |
| Валидация DTO | `go-playground/validator/v10` |
| JWT | `golang-jwt/jwt/v5` (HMAC, алгоритм пиним) |
| Хэш паролей | `golang-jwt/jwt/v5` (HMAC) |
| Redis | `redis/go-redis/v9` |
| Хранилище | `minio/minio-go/v7` |
| Логи | `go.uber.org/zap` (structured JSON) |
| Тесты | `stretchr/testify` + `testcontainers-go` (postgres/redis) |

---

## 0. Каркас проекта
- [x] `go mod init`, структура пакетов выше
- [x] `internal/config` — структуры + `Load()` (godotenv+env) + `Validate()`, разбит по файлам
- [x] `pkg/logger` — zap factory (конфигурируется из `config.logger`)
- [x] `adapter/database/postgres/client.go` — пул pgx + ping
- [x] `adapter/cache/redis/client.go` — клиент go-redis + ping
- [x] `pkg/httpx` — JSON/error-хелперы, request_id + claims в context, extract Bearer
- [x] `internal/app/container.go` (DI: конфиг→клиенты→репо→кэш→сервисы→хендлеры) и `app.go` (запуск + graceful shutdown по `SERVER_SHUTDOWN_TIMEOUT`)
- [x] `cmd/api/main.go`, `GET /health`
- [x] `Dockerfile` (multi-stage, production-ready), `docker-compose.yml` (api + postgres + redis + rustfs), `.env.example`

## 1. БД, миграции, sqlc
- [x] `sqlc.yaml` (engine postgresql, sql_package pgx/v5)
- [x] Миграции goose: `users`, `subscriptions`, `artists`, `albums`, `genres`, `tracks`, `playlists`, `playlist_tracks`, `favorites`, `listening_history`
- [x] Seed-миграция: базовый список `genres`, `artists`, `albums`, `tracks` для демо
- [x] Индексы: поиск треков, FK, уникальность (`favorites(user_id,track_id)`, `artists(name)`, `albums(name)` и т.п.)
- [x] `database/queries/*.sql` под каждый домен → `sqlc generate`
- [x] Миграции прогоняются автоматически при старте (goose)

## 2. Auth (`service/auth`, `v1/auth`)
- [x] `domain/user/password.go` — bcrypt + salt (`NewPassword`/`Matches`)
- [x] `pkg/jwt` manager: подпись/парсинг HMAC, пин алгоритма, проверка `exp`, claims (`user_id`, `role`)
- [x] Refresh-токены — JWT + Redis allowlist без таблицы `refresh_tokens`
- [x] `cache/blacklist` — отзыв access-токена при logout (Redis, TTL = остаток жизни токена)
- [x] `POST /auth/register` (роль USER, подписка FREE по умолчанию)
- [x] `POST /auth/login` (access + refresh)
- [x] `POST /auth/refresh` (ротация refresh)
- [x] `POST /auth/logout` (blacklist access; remove refresh из allowlist)

## 3. Middleware
- [x] `adapter/.../middleware/request_id` (correlation id в context)
- [x] `adapter/.../middleware/logger` (метод, путь, статус, длительность, request_id, user_id)
- [x] `adapter/.../middleware/recover` (восстановление после panic)
- [x] `adapter/.../middleware/auth` (Bearer → валидация подписи+exp → проверка blacklist → identity в context)
- [x] `adapter/.../middleware/role` (`RequireRole("ADMIN")` поверх auth)


## 4. Users (`service/user`, `v1/user`)
- [x] `GET /users/me` (id, email, username, role, subscription_type, created_at)
- [x] `PATCH /users/me`

## 5. Tracks (`service/track`, `v1/track`)
- [x] `GET /tracks` (пагинация `page`/`limit`)
- [x] `GET /tracks/{id}`
- [x] `GET /tracks/search` (по названию/исполнителю/жанру/альбому, JOIN, без N+1)
- [x] `POST /tracks/{id}/play` (существует → запись в `listening_history` → возврат трека)

## 6. Playlists (`service/playlist`, `v1/playlist`) — только свои
- [x] `POST /playlists` (FREE: лимит из `FREE_PLAYLIST_LIMIT`, проверка в сервисе)
- [x] `GET /playlists`
- [x] `GET /playlists/{id}` (проверка владельца)
- [x] `PUT /playlists/{id}`
- [x] `DELETE /playlists/{id}`
- [x] `POST /playlists/{playlist_id}/tracks/{track_id}`
- [x] `DELETE /playlists/{playlist_id}/tracks/{track_id}`

## 7. Favorites (`service/favorite`, `v1/favorite`)
- [x] `POST /favorites/tracks/{track_id}` (FREE: лимит из `FREE_FAVORITES_LIMIT`, проверка в сервисе)
- [x] `GET /favorites/tracks`
- [x] `DELETE /favorites/tracks/{track_id}`

## 8. Listening history (`service/history`, `v1/history`)
- [x] `GET /listening-history` (id трека, название, исполнитель, listened_at)

## 9. Admin (`v1/admin`, под `role` middleware)
- [x] `POST /admin/tracks` (multipart/form-data: загрузка в S3, создание в БД)
- [x] `PUT /admin/tracks/{id}` (инвалидация `track:{id}`)
- [x] `DELETE /admin/tracks/{id}` (soft delete: `deleted_at=NOW()`; инвалидация `track:{id}`)
- [x] `PATCH /admin/users/{id}/subscription`

## 10. Кэш (`cache/*`, Redis)
- [x] `track:{id}` — трек по id (инвалидация при update/delete админом)
- [x] `genres` — список жанров (сид-справочник)
- [x] результаты поиска
- [x] cache-aside: read-through + инвалидация при записи

## 11. Обработка ошибок
- [x] Доменные ошибки (`domain/<d>/errors.go`) → маппинг в HTTP в хендлере
- [x] JSON-формат `{ "error": "..." }`
- [x] Коды: 200, 201, 400, 401, 403, 404, 409, 500

## 12. Swagger / OpenAPI
- [x] `api/v1` (openapi.yaml + paths + components), securityScheme Bearer
- [x] Раздача через `web/swagger/index.html` + `swagger/routes.go`

## 13. Тесты
- [x] Сервисы: unit + `mocks_test.go` (intf потребителя), require/assert, error-пути
- [x] Хендлеры: `httptest` (статус, JSON, валидация)
- [x] Обязательно: register, login, profile, создание плейлиста, добавление в избранное, лимит плейлистов FREE, доступ к admin-эндпоинтам (готовность 100%)

## 14. README
- [x] Инструкция запуска (Docker Compose + локально) — проверяется на защите

## 15. Готовность к защите
- [x] Старт через `docker-compose up`
- [x] Объяснить слои (handler/service/repository), DI, поток `handler → service → repository`
- [x] Объяснить схему БД и sqlc/goose
- [x] Показать JWT-флоу (login → access/refresh → logout/blacklist)
- [x] Показать API через Swagger

---

## Опционально (для усиления)
- [x] Загрузка обложек альбомов (S3/RustFS) — реализовано для треков
- [x] Ротация refresh-токенов
- [x] Graceful shutdown (если не сделан в §0)
