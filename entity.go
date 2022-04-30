package kumparan_assessment

import "time"

// Article .
type Article struct {
	ID      int64     `json:"id"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

// ArticleFilter .
type ArticleFilter struct {
	Author string
	Search string
}
