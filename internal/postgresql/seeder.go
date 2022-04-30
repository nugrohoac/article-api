package postgresql

import (
	"database/sql"
	"testing"
	
	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	ka "github.com/nugrohoac/kumparan-assessment"
)

// SeedArticles .
func SeedArticles(t *testing.T, db *sql.DB, articles ...ka.Article) {
	qInsert := sq.Insert("article").
		Columns("author",
			"title",
			"body",
		)

	for _, article := range articles {
		qInsert = qInsert.Values(article.Author, article.Title, article.Body)
	}

	query, args, err := qInsert.PlaceholderFormat(sq.Dollar).ToSql()
	require.NoError(t, err)

	_, err = db.Exec(query, args...)
	require.NoError(t, err)
}
