package cache

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"

	ka "github.com/nugrohoac/kumparan-assessment"
)

// KeyArticle .
const KeyArticle = "article"

type articleRedis struct {
	rdb        *redis.Client
	repository ka.ArticleRepository
}

// Store .
func (a articleRedis) Store(ctx context.Context, article ka.Article) (ka.Article, error) {
	article, err := a.repository.Store(ctx, article)
	if err != nil {
		return ka.Article{}, err
	}

	if err = a.rdb.Del(ctx, KeyArticle).Err(); err != nil {
		logrus.Error("failed delete article")
	}

	return article, nil
}

// Fetch .
func (a articleRedis) Fetch(ctx context.Context, filter ka.ArticleFilter) ([]ka.Article, error) {
	articles := make([]ka.Article, 0)

	value, err := a.rdb.Get(ctx, KeyArticle).Result()
	if err != nil {
		if err == redis.Nil {
			articles, err = a.repository.Fetch(ctx, filter)
			if err != nil {
				return nil, err
			}

			byteArticles, err := json.Marshal(articles)
			if err != nil {
				return nil, err
			}

			if err = a.rdb.Set(ctx, KeyArticle, string(byteArticles), 0).Err(); err != nil {
				return nil, err
			}

			return articles, nil
		}

		return nil, err
	}

	if err = json.Unmarshal([]byte(value), &articles); err != nil {
		return nil, err
	}

	return articles, nil
}

// NewArticleRedis .
func NewArticleRedis(rdb *redis.Client, repository ka.ArticleRepository) ka.ArticleRepository {
	return articleRedis{
		rdb:        rdb,
		repository: repository,
	}
}
