# 🎵 Music Platform API

Backend-сервис музыкальной платформы.

## 🛠 Технологический стек

| Категория | Технология |
|---|---|
| Язык | Go (1.26+) |
| Архитектура | Чистая (Clean Architecture): Handler → Service → Repository |
| БД | PostgreSQL 18 (SQL-миграции goose, запросы sqlc) |
| Кэш | Redis 8.0 |
| Хранилище | MinIO (S3-совместимое) |
| Аутентификация | JWT (Access + Refresh) |
| Логирование | zap (структурированное) |
| Документация | Swagger (OpenAPI 3.1) |
| Контейнеризация | Docker + Docker Compose |

## 🚀 Быстрый старт

### Требования
* Docker и Docker Compose

### Запуск через Docker
1. Скопируйте пример конфига:
   ```bash
   cp .env.example .env
   ```
2. Запустите инфраструктуру:
   ```bash
   docker-compose up -d --build
   ```
3. Наполните S3 тестовыми данными (треками):
   ```bash
   ./scripts/seed-s3.sh
   ```

API будет доступно по адресу `http://localhost:8080`.
Миграции БД (включая seed-данные) применяются автоматически при старте.

## 📖 Документация

Swagger UI доступен по адресу: `http://localhost:8080/swagger`

Основные эндпоинты:
* `POST /api/v1/auth/register` — Регистрация
* `POST /api/v1/auth/login` — Вход
* `GET /api/v1/tracks` — Каталог треков
* `POST /api/v1/admin/tracks` — Загрузка нового трека (multipart/form-data)

## 🏗 Архитектура

Проект построен по принципам Clean Architecture:
* `internal/domain` — Бизнес-модели и интерфейсы
* `internal/service` — Бизнес-логика и координация слоев
* `internal/repository` — Реализация работы с БД (PostgreSQL)
* `internal/adapter/transport/http` — Транспортный слой (chi router)
* `internal/adapter/storage` — Работа с внешними хранилищами (S3)
* `internal/cache` — Слой кэширования (Redis)
