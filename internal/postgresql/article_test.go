package postgresql_test

import (
	"context"
	ka "github.com/nugrohoac/kumparan-assessment"
	"github.com/nugrohoac/kumparan-assessment/testdata"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/nugrohoac/kumparan-assessment/internal/postgresql"
)

type articleSuite struct {
	postgresql.TestSuite
}

func TestArticleRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test article")
	}

	suite.Run(t, new(articleSuite))
}

func (a articleSuite) TestArticleRepo_Store() {
	articles := make([]ka.Article, 0)
	testdata.GoldenJSONUnmarshal(a.T(), "articles", &articles)

	articleRepository := postgresql.NewArticleRepository(a.DBConn)
	article, err := articleRepository.Store(context.Background(), articles[0])

	require.NoError(a.T(), err)
	require.Equal(a.T(), article.Title, articles[0].Title)
	require.Equal(a.T(), article.Body, articles[0].Body)
	require.Equal(a.T(), article.Author, articles[0].Author)
}

func (a articleSuite) TestArticleRepo_Fetch() {
	articles := make([]ka.Article, 0)
	testdata.GoldenJSONUnmarshal(a.T(), "articles", &articles)

	postgresql.SeedArticles(a.T(), a.DBConn, articles...)

	articleRepository := postgresql.NewArticleRepository(a.DBConn)
	response, err := articleRepository.Fetch(context.Background(), ka.ArticleFilter{})
	require.NoError(a.T(), err)
	require.Equal(a.T(), 3, len(response))

	response, err = articleRepository.Fetch(context.Background(), ka.ArticleFilter{Author: "JHON"})
	require.NoError(a.T(), err)
	require.Equal(a.T(), 1, len(response))

	response, err = articleRepository.Fetch(context.Background(), ka.ArticleFilter{Search: "second"})
	require.NoError(a.T(), err)
	require.Equal(a.T(), 2, len(response))
}
