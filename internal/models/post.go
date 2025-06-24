package models

import (
	"time"
)

type Post struct {
    ID        int       `db:"id" json:"id"`
    Title     string    `db:"title" json:"title"`
    Content   string    `db:"content" json:"content"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    Comments  []*Comment `db:"-"`
}

type Comment struct {
    ID        int       `db:"id"`
    PostID    int       `db:"post_id"`
    Content   string    `db:"content"`
    CreatedAt time.Time `db:"created_at"`
}
