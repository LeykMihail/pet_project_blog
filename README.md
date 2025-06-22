# Pet Project Blog

Простой блог на Go с использованием Gin, чистой архитектуры и миграций.

## Структура проекта

```
pet_project_blog/
  cmd/app/main.go         # Точка входа
  internal/
    config/               # Конфигурация
    handlers/             # HTTP-обработчики
    migrations/           # Миграции БД
    models/               # Модели данных
    repository/           # Работа с БД
    services/             # Бизнес-логика
```

## Основные возможности
- Просмотр всех постов (`GET /posts`)
- Получение поста по ID (`GET /posts/:id`)
- Создание поста (`POST /posts`)
- Гибкая фильтрация полей через query-параметр `fields`

## Запуск

1. Установите зависимости:
   ```sh
   go mod download
   ```
2. Примените миграции к вашей базе данных (например, с помощью [golang-migrate](https://github.com/golang-migrate/migrate))
3. Запустите приложение:
   ```sh
   go run cmd/app/main.go
   ```

## Пример запросов

- Получить все посты:
  ```sh
  curl http://localhost:8080/posts
  ```
- Получить только id и title:
  ```sh
  curl http://localhost:8080/posts?fields=id,title
  ```
- Получить пост по ID:
  ```sh
  curl http://localhost:8080/posts/1
  ```
- Создать пост:
  ```sh
  curl -X POST http://localhost:8080/posts -H 'Content-Type: application/json' -d '{"title":"Заголовок","content":"Текст поста"}'
  ```
