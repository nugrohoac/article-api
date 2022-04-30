package cache_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	ka "github.com/nugrohoac/kumparan-assessment"
	"github.com/nugrohoac/kumparan-assessment/internal/cache"
	"github.com/nugrohoac/kumparan-assessment/mocks"
	"github.com/nugrohoac/kumparan-assessment/testdata"
)

type articleSuite struct {
	cache.TestSuite
}

func TestArticleRedis(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test article")
	}

	suite.Run(t, new(articleSuite))
}

func (a articleSuite) TestArticleRedis_Store() {
	articles := make([]ka.Article, 0)
	testdata.GoldenJSONUnmarshal(a.T(), "articles", &articles)

	article := articles[0]
	articleRepoMock := new(mocks.ArticleRepository)
	articleRedis := cache.NewArticleRedis(a.Client, articleRepoMock)

	byteArticle, err := json.Marshal(article)
	require.NoError(a.T(), err)
	err = a.Client.Set(context.Background(), cache.KeyArticle, string(byteArticle), 0).Err()
	require.NoError(a.T(), err)

	articleRepoMock.On("Store", mock.Anything, article).Return(article, nil).Once()
	response, err := articleRedis.Store(context.Background(), article)
	require.NoError(a.T(), err)
	require.Equal(a.T(), response, article)

	err = a.Client.Get(context.Background(), cache.KeyArticle).Err()
	require.Equal(a.T(), redis.Nil, err)

	articleRepoMock.AssertExpectations(a.T())
}

func (a articleSuite) TestArticleRedis_Fetch() {
	articles := make([]ka.Article, 0)
	testdata.GoldenJSONUnmarshal(a.T(), "articles", &articles)

	articleRepoMock := new(mocks.ArticleRepository)
	articleRedis := cache.NewArticleRedis(a.Client, articleRepoMock)

	articleRepoMock.On("Fetch", mock.Anything, ka.ArticleFilter{}).Return(articles, nil).Once()
	response, err := articleRedis.Fetch(context.Background(), ka.ArticleFilter{})
	require.NoError(a.T(), err)
	require.Equal(a.T(), articles, response)

	value, err := a.Client.Get(context.Background(), cache.KeyArticle).Result()
	require.NoError(a.T(), err)

	_articles := make([]ka.Article, 0)
	err = json.Unmarshal([]byte(value), &_articles)
	require.NoError(a.T(), err)
	require.Equal(a.T(), articles, _articles)

	articleRepoMock.AssertExpectations(a.T())
}
