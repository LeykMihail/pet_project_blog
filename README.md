# Pet Project Blog

Простой блог на Go с использованием Gin, чистой архитектуры и миграций PostgreSQL.

## 📦 Структура проекта

```
pet_project_blog/
  cmd/app/main.go         # Точка входа
  internal/
    config/               # Конфигурация и переменные окружения
    handlers/             # HTTP-обработчики (Gin)
    migrations/           # Миграции БД (SQL)
    models/               # Модели данных
    repository/           # Работа с БД (PostgreSQL)
    services/             # Бизнес-логика
```

## 🚀 Быстрый старт

1. **Клонируйте репозиторий:**
   ```sh
   git clone <your-repo-url>
   cd pet_project_blog
   ```
2. **Установите зависимости:**
   ```sh
   go mod download
   ```
3. **Настройте переменные окружения:**
   - `APP_PORT` — порт HTTP-сервера (по умолчанию 8080)
   - `DB_CONN_STR` — строка подключения к PostgreSQL
     
     Пример:
     ```sh
     export APP_PORT=8080
     export DB_CONN_STR="postgres://user:password@localhost:5432/blog?sslmode=disable"
     ```
4. **Примените миграции:**
   ```sh
   # Установите golang-migrate, если не установлен
   migrate -path internal/migrations -database "$DB_CONN_STR" up
   ```
5. **Запустите приложение:**
   ```sh
   go run cmd/app/main.go
   ```

## 🛠️ Возможности
- Просмотр всех постов (`GET /posts`)
- Получение поста по ID с комментариями (`GET /posts/:id`)
- Создание поста (`POST /posts`)
- Добавление комментария к посту (`POST /posts/:id/comments`)
- Получение комментариев поста (`GET /posts/:id/comments`)
- Гибкая фильтрация полей через query-параметр `fields`
- Регистрация пользователя (`POST /register`)
- Вход пользователя (логин) (`POST /login`)
- Авторизация через cookie (требуется для создания постов и комментариев)

## 🧪 Примеры запросов

- Получить все посты:
  ```sh
  curl http://localhost:8080/posts
  ```
- Получить только id и title:
  ```sh
  curl http://localhost:8080/posts?fields=id,title
  ```
- Получить пост по ID с комментариями:
  ```sh
  curl http://localhost:8080/posts/1
  ```
- Создать пост:
  ```sh
  curl -X POST http://localhost:8080/posts -H 'Content-Type: application/json' -d '{"title":"Заголовок","content":"Текст поста"}'
  ```
- Добавить комментарий:
  ```sh
  curl -X POST http://localhost:8080/posts/1/comments -H 'Content-Type: application/json' -d '{"content":"Комментарий"}'
  ```
- Получить комментарии к посту:
  ```sh
  curl http://localhost:8080/posts/1/comments
  ```
- Зарегистрировать пользователя:
  ```sh
  curl -X POST http://localhost:8080/register -H 'Content-Type: application/json' -d '{"email":"user@example.com","password":"password123"}'
  ```
- Войти (логин):
  ```sh
  curl -X POST http://localhost:8080/login -H 'Content-Type: application/json' -d '{"email":"user@example.com","password":"password123"}'
  ```

- Получить JWT-токен:
  После успешного логина в ответе будет поле `token`. Скопируйте его для дальнейших запросов.

- Использовать токен для авторизованных запросов (например, создать пост):
  ```sh
  curl -X POST http://localhost:8080/posts \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer <ваш_токен>' \
    -d '{"title":"Заголовок","content":"Текст поста"}'
  ```
