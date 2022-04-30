package kumparan_assessment

import "time"

// Article .
type Article struct {
	ID      int64     `json:"id"`
	Author  string    `json:"author" validate:"required"`
	Title   string    `json:"title" validate:"required"`
	Body    string    `json:"body" validate:"required"`
	Created time.Time `json:"created"`
}

// ArticleFilter .
type ArticleFilter struct {
	Author string
	Search string
}
