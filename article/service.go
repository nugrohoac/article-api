package article

import (
	"context"

	ka "github.com/nugrohoac/kumparan-assessment"
)

type articleService struct {
	articleRepository ka.ArticleRepository
}

// Store .
func (a articleService) Store(ctx context.Context, article ka.Article) (ka.Article, error) {
	return a.articleRepository.Store(ctx, article)
}

// Fetch .
func (a articleService) Fetch(ctx context.Context, filter ka.ArticleFilter) ([]ka.Article, error) {
	return a.articleRepository.Fetch(ctx, filter)
}

// NewArticleService .
func NewArticleService(articleRepository ka.ArticleRepository) ka.ArticleService {
	return articleService{articleRepository: articleRepository}
}
