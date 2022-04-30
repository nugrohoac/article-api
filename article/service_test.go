package article_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ka "github.com/nugrohoac/kumparan-assessment"
	_article "github.com/nugrohoac/kumparan-assessment/article"
	"github.com/nugrohoac/kumparan-assessment/mocks"
	"github.com/nugrohoac/kumparan-assessment/testdata"
)

func TestArticleService_Store(t *testing.T) {
	var articles []ka.Article
	testdata.GoldenJSONUnmarshal(t, "articles", &articles)
	article := articles[0]

	tests := map[string]struct {
		paramArticle ka.Article
		storeArticle testdata.FuncCaller
		expectedResp ka.Article
		expectedErr  error
	}{
		"error": {
			paramArticle: article,
			storeArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, article},
				Output:   []interface{}{ka.Article{}, errors.New("error")},
			},
			expectedResp: ka.Article{},
			expectedErr:  errors.New("error"),
		},
		"success": {
			paramArticle: article,
			storeArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, article},
				Output:   []interface{}{article, nil},
			},
			expectedResp: article,
			expectedErr:  nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			articleRepoMock := new(mocks.ArticleRepository)

			if test.storeArticle.IsCalled {
				articleRepoMock.On("Store", test.storeArticle.Input...).
					Return(test.storeArticle.Output...).
					Once()
			}

			articleService := _article.NewArticleService(articleRepoMock)
			response, err := articleService.Store(context.Background(), test.paramArticle)

			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.expectedResp, response)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedResp, response)
		})
	}
}

func TestArticleService_Fetch(t *testing.T) {
	var articles []ka.Article
	testdata.GoldenJSONUnmarshal(t, "articles", &articles)

	tests := map[string]struct {
		paramFilter  ka.ArticleFilter
		fetchArticle testdata.FuncCaller
		expectedResp []ka.Article
		expectedErr  error
	}{
		"error": {
			paramFilter: ka.ArticleFilter{},
			fetchArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, ka.ArticleFilter{}},
				Output:   []interface{}{nil, errors.New("error")},
			},
			expectedResp: nil,
			expectedErr:  errors.New("error"),
		},
		"success": {
			paramFilter: ka.ArticleFilter{},
			fetchArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, ka.ArticleFilter{}},
				Output:   []interface{}{articles, nil},
			},
			expectedResp: articles,
			expectedErr:  nil,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			articleRepoMock := new(mocks.ArticleRepository)

			if test.fetchArticle.IsCalled {
				articleRepoMock.On("Fetch", test.fetchArticle.Input...).
					Return(test.fetchArticle.Output...).
					Once()
			}

			articleService := _article.NewArticleService(articleRepoMock)
			response, err := articleService.Fetch(context.Background(), test.paramFilter)

			if err != nil {
				require.Error(t, err)
				require.Equal(t, test.expectedResp, response)

				return
			}

			require.NoError(t, err)
			require.Equal(t, test.expectedResp, response)
		})
	}
}
