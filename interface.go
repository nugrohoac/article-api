package kumparan_assessment

import "context"

// ArticleRepository .
type ArticleRepository interface {
	Store(ctx context.Context, article Article) (Article, error)
	Fetch(ctx context.Context, filter ArticleFilter) ([]Article, error)
}

// ArticleService .
type ArticleService interface {
	Store(ctx context.Context, article Article) (Article, error)
	Fetch(ctx context.Context, filter ArticleFilter) ([]Article, error)
}
