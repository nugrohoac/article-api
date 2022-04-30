package postgresql

import (
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	ka "github.com/nugrohoac/kumparan-assessment"
)

type articleRepo struct {
	db *sql.DB
}

// Store .
func (a articleRepo) Store(ctx context.Context, article ka.Article) (ka.Article, error) {
	article.Created = time.Now()

	query, args, err := sq.Insert("article").Columns("author",
		"title",
		"body",
		"created",
	).Values(article.Author,
		article.Title,
		article.Body,
		article.Created,
	).PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return ka.Article{}, err
	}

	if _, err = a.db.ExecContext(ctx, query, args...); err != nil {
		return ka.Article{}, err
	}

	return article, nil
}

// Fetch .
func (a articleRepo) Fetch(ctx context.Context, filter ka.ArticleFilter) ([]ka.Article, error) {
	qSelect := sq.Select("id",
		"author",
		"title",
		"body",
		"created",
	).From("article").
		OrderBy("created desc").
		PlaceholderFormat(sq.Dollar)

	if filter.Author != "" {
		qSelect = qSelect.Where(sq.Eq{"LOWER(author)": strings.ToLower(filter.Author)})
	}

	if filter.Search != "" {
		search := "%" + strings.ToLower(filter.Search) + "%"

		qSelect = qSelect.Where(sq.Or{
			sq.Like{"LOWER(title)": search},
			sq.Like{"LOWER(body)": search},
		})
	}

	query, args, err := qSelect.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if errClose := rows.Close(); errClose != nil {
			logrus.Error(errClose)
		}
	}()

	articles := make([]ka.Article, 0)

	for rows.Next() {
		var article ka.Article

		if err = rows.Scan(
			&article.ID,
			&article.Author,
			&article.Title,
			&article.Body,
			&article.Created,
		); err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

// NewArticleRepository .
func NewArticleRepository(db *sql.DB) ka.ArticleRepository {
	return articleRepo{db: db}
}
