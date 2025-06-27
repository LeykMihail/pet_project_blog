package models

import (
	"time"
)

type Post struct {
    ID        int       `db:"id"`
    Title     string    `db:"title" `
    Content   string    `db:"content"`
    CreatedAt time.Time `db:"created_at"`
    UserID    int       `db:"user_id"`
    Comments  []*Comment `db:"-"`
}

type Comment struct {
    ID        int       `db:"id"`
    PostID    int       `db:"post_id"`
    Content   string    `db:"content"`
    CreatedAt time.Time `db:"created_at"`
    UserID    int       `db:"user_id"`
}

type User struct {
    ID           int       `db:"id"`
    Email        string    `db:"email"`
    PasswordHash string    `db:"password_hash"`
    CreatedAt    time.Time `db:"created_at"`
}
