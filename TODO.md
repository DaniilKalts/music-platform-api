# TODO

Порядок сборки — сверху вниз. Отмечай пункты по мере готовности.

**Архитектура — как в `rbk-school/7-week/api-service`** (Clean Architecture, один сервис). Конвенции, нейминг пакетов и тулчейн повторяем оттуда. **Ничего сверх ТЗ не добавляем** (никакого gateway, лишних доменов, Kafka и т.п.).

> **Soft delete:** используем `is_active` (boolean) в `tracks`. `is_deleted` из ТЗ игнорируем.
>
> **Аудиофайлы:** лежат в бакете RustFS (S3). Бэкенд **не работает с S3 в коде** — `tracks.file_url` это обычная строка-ссылка на объект. Файл заливается в бакет вне API, админ передаёт готовый URL в `POST /admin/tracks`. Никакого S3 SDK и upload-эндпоинта.
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
| Хэш паролей | `golang.org/x/crypto/bcrypt` |
| Redis | `redis/go-redis/v9` |
| Логи | `go.uber.org/zap` (structured JSON) |
| Тесты | `stretchr/testify` + `testcontainers-go` (postgres/redis) |

---

## Целевая структура пакетов

```
cmd/api/main.go
internal/
  app/            # app.go (запуск + graceful shutdown), container.go (DI)
  config/         # config.go, load.go, server.go, postgres.go, redis.go, jwt.go, logger.go
  domain/
    user/         # model.go, errors.go, password.go, role.go, subscription.go
    track/        # model.go, errors.go  (+ artist/album/genre как поля/модели каталога)
    playlist/     # model.go, errors.go
    favorite/     # model.go, errors.go
    history/      # model.go, errors.go
  service/        # services.go + <domain>/{service.go, service_test.go, mocks_test.go}
    auth/ user/ track/ playlist/ favorite/ history/
  repository/     # repositories.go + <domain>/{repository.go, converter.go, repository_integration_test.go}
    user/ track/ playlist/ favorite/ history/
  cache/          # caches.go + <name>/{cache.go, cache_integration_test.go}
    track/ genre/ popular/ search/ blacklist/
  adapter/
    database/postgres/   # client.go, errors.go, sqlc/
    cache/redis/         # client.go
    transport/http/
      router.go
      middleware/        # request_id.go, logger.go, recover.go, auth.go, role.go
      swagger/routes.go
      v1/                # routes.go + <domain>/{handler.go, dto.go, routes.go, handler_test.go, mocks_test.go}
pkg/
  logger/         # zap factory
  jwt/            # JWT manager (подпись/парсинг, claims)
  httpx/          # JSON/error-хелперы, request_id + claims в context, extract Bearer
database/
  migrations/     # goose
  queries/        # sqlc-исходники *.sql
api/v1/           # OpenAPI: openapi.yaml, paths/, components/{schemas,examples,...}
web/swagger/index.html
sqlc.yaml  Dockerfile  docker-compose.yml  .env.example  go.mod
```

> **Definition of Done на домен:** миграция (goose) → запрос (`database/queries/*.sql`) → `sqlc generate` → `repository/<d>` (`repository.go` + `converter.go` sqlc-модель↔домен) → `service/<d>/service.go` → `adapter/.../v1/<d>` (`handler.go` + `dto.go` домен↔DTO + `routes.go`) → тесты (service+mocks, handler+httptest, repository+testcontainers).

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
- [ ] `Dockerfile` (multi-stage, non-root), `docker-compose.yml` (api + postgres + redis + rustfs), `.env.example`

**Env:** `APP_PORT`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `REDIS_HOST`, `REDIS_PORT`, `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `FREE_PLAYLIST_LIMIT`, `FREE_FAVORITES_LIMIT`

## 1. БД, миграции, sqlc
- [x] `sqlc.yaml` (engine postgresql, sql_package pgx/v5)
- [ ] Миграции goose: `users`, `subscriptions`, `artists`, `albums`, `genres`, `tracks`, `playlists`, `playlist_tracks`, `favorites`, `listening_history`
- [ ] Seed-миграция: фиксированный список `genres` (справочник, только чтение)
- [ ] Индексы: поиск треков, FK, уникальность (`favorites(user_id,track_id)`, `artists(name)`, `albums(name)` и т.п.)
- [ ] `database/queries/*.sql` под каждый домен → `sqlc generate`
- [x] Миграции прогоняются автоматически при старте (goose)

## 2. Auth (`service/auth`, `v1/auth`)
- [x] `domain/user/password.go` — bcrypt + salt (`NewPassword`/`Matches`)
- [x] `pkg/jwt` manager: подпись/парсинг HMAC, пин алгоритма, проверка `exp`, claims (`user_id`, `role`)
- [x] Refresh-токены — stateless JWT без таблицы `refresh_tokens`; хранение добавлять только если понадобится server-side revoke/rotation
- [x] `cache/blacklist` — отзыв access-токена при logout (Redis, TTL = остаток жизни токена)
- [x] `POST /auth/register` (роль USER, подписка FREE по умолчанию)
- [x] `POST /auth/login` (access + refresh)
- [x] `POST /auth/refresh` (ротация refresh)
- [x] `POST /auth/logout` (blacklist access; refresh stateless)

## 3. Middleware
- [x] `adapter/.../middleware/request_id` (correlation id в context)
- [x] `adapter/.../middleware/logger` (метод, путь, статус, длительность, request_id, user_id)
- [x] `adapter/.../middleware/recover` (восстановление после panic)
- [x] `adapter/.../middleware/auth` (Bearer → валидация подписи+exp → проверка blacklist → identity в context)
- [x] `adapter/.../middleware/role` (`RequireRole("ADMIN")` поверх auth)


## 4. Users (`service/user`, `v1/user`)
- [ ] `GET /users/me` (id, email, username, role, subscription_type, created_at)
- [ ] `PATCH /users/me`

## 5. Tracks (`service/track`, `v1/track`)
- [ ] `GET /tracks` (пагинация `page`/`limit`)
- [ ] `GET /tracks/{id}`
- [ ] `GET /tracks/search` (по названию/исполнителю/жанру/альбому, JOIN, без N+1)
- [ ] `POST /tracks/{id}/play` (существует → запись в `listening_history` → возврат трека)

## 6. Playlists (`service/playlist`, `v1/playlist`) — только свои
- [ ] `POST /playlists` (FREE: лимит из `FREE_PLAYLIST_LIMIT`, проверка в сервисе)
- [ ] `GET /playlists`
- [ ] `GET /playlists/{id}` (проверка владельца)
- [ ] `PUT /playlists/{id}`
- [ ] `DELETE /playlists/{id}`
- [ ] `POST /playlists/{playlist_id}/tracks/{track_id}`
- [ ] `DELETE /playlists/{playlist_id}/tracks/{track_id}`

## 7. Favorites (`service/favorite`, `v1/favorite`)
- [ ] `POST /favorites/tracks/{track_id}` (FREE: лимит из `FREE_FAVORITES_LIMIT`, проверка в сервисе)
- [ ] `GET /favorites/tracks`
- [ ] `DELETE /favorites/tracks/{track_id}`

## 8. Listening history (`service/history`, `v1/history`)
- [ ] `GET /listening-history` (id трека, название, исполнитель, listened_at)

## 9. Admin (`v1/admin`, под `role` middleware)
- [ ] `POST /admin/tracks` (артист/альбом — find-or-create по имени, жанр — по сид-справочнику; всё в одной транзакции)
- [ ] `PUT /admin/tracks/{id}` (тот же find-or-create; инвалидация `track:{id}`)
- [ ] `DELETE /admin/tracks/{id}` (soft delete: `is_active=false`; инвалидация `track:{id}`)
- [ ] `PATCH /admin/users/{id}/subscription`

## 10. Кэш (`cache/*`, Redis)
- [ ] `track:{id}` — трек по id (инвалидация при update/delete админом)
- [ ] `genres` — список жанров (сид-справочник)
- [ ] результаты поиска
- [ ] cache-aside: read-through + инвалидация при записи
- [ ] `popular_tracks` — только вместе с опциональной фичей «топ популярных» (§Опционально)

## 11. Обработка ошибок
- [x] Доменные ошибки (`domain/<d>/errors.go`) → маппинг в HTTP в хендлере
- [x] JSON-формат `{ "error": "..." }`
- [x] Коды: 200, 201, 400, 401, 403, 404, 409, 500

## 12. Swagger / OpenAPI
- [ ] `api/v1` (openapi.yaml + paths + components), securityScheme Bearer
- [ ] Раздача через `web/swagger/index.html` + `swagger/routes.go`

## 13. Тесты
- [x] Сервисы: unit + `mocks_test.go` (intf потребителя), require/assert, error-пути
- [x] Хендлеры: `httptest` (статус, JSON, валидация)
- [ ] Репозитории: `repository_integration_test.go` (testcontainers postgres, `-tags=integration`)
- [ ] Кэш: integration (testcontainers redis, `-tags=integration`)
- [x] Обязательно: register, login, profile, создание плейлиста, добавление в избранное, лимит плейлистов FREE, доступ к admin-эндпоинтам (частично готово: register, login)

## 14. README
- [ ] Инструкция запуска (Docker Compose + локально) — проверяется на защите

## 15. Готовность к защите
- [ ] Старт через `docker-compose up`
- [ ] Объяснить слои (handler/service/repository), DI, поток `handler → service → repository`
- [ ] Объяснить схему БД и sqlc/goose
- [ ] Показать JWT-флоу (login → access/refresh → logout/blacklist)
- [ ] Показать API через Swagger

---

## Опционально (для усиления)
- [ ] Топ популярных треков
- [ ] Рекомендации на основе истории
- [ ] Загрузка обложек альбомов (S3/RustFS)
- [ ] Rate limiting
- [ ] Ротация refresh-токенов
- [ ] Audit log действий администратора
- [ ] Graceful shutdown (если не сделан в §0)
- [ ] Prometheus metrics
