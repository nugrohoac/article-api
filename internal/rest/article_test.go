package rest_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ka "github.com/nugrohoac/kumparan-assessment"
	"github.com/nugrohoac/kumparan-assessment/internal/rest"
	"github.com/nugrohoac/kumparan-assessment/mocks"
	"github.com/nugrohoac/kumparan-assessment/testdata"
)

func TestRegisterPathArticle_Store(t *testing.T) {
	e := echo.New()
	v := validator.New()

	var articles []ka.Article
	testdata.GoldenJSONUnmarshal(t, "articles", &articles)
	article := articles[0]

	byteArticle, err := json.Marshal(article)
	require.NoError(t, err)

	tests := map[string]struct {
		body               []byte
		path               string
		storeArticle       testdata.FuncCaller
		expectedResp       ka.Article
		expectedStatusCode int
	}{
		"bad request": {
			body:               nil,
			path:               "/articles",
			storeArticle:       testdata.FuncCaller{},
			expectedResp:       ka.Article{},
			expectedStatusCode: http.StatusBadRequest,
		},
		"internal server error": {
			body: byteArticle,
			path: "/articles",
			storeArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, article},
				Output:   []interface{}{ka.Article{}, errors.New("error")},
			},
			expectedResp:       ka.Article{},
			expectedStatusCode: http.StatusInternalServerError,
		},
		"success": {
			body: byteArticle,
			path: "/articles",
			storeArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, article},
				Output:   []interface{}{article, nil},
			},
			expectedResp:       article,
			expectedStatusCode: http.StatusCreated,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			articleServiceMock := new(mocks.ArticleService)

			if test.storeArticle.IsCalled {
				articleServiceMock.On("Store", test.storeArticle.Input...).
					Return(test.storeArticle.Output...).
					Once()
			}

			req := httptest.NewRequest(http.MethodPost, test.path, strings.NewReader(string(test.body)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			rest.RegisterPathArticle(e, v, articleServiceMock)
			e.ServeHTTP(rec, req)

			require.Equal(t, test.expectedStatusCode, rec.Code)

			if rec.Code == http.StatusCreated {
				var _article ka.Article
				err = json.Unmarshal(rec.Body.Bytes(), &_article)
				require.NoError(t, err)

				require.Equal(t, test.expectedResp, _article)
			}

			articleServiceMock.AssertExpectations(t)
		})
	}
}

func TestRegisterPathArticle_Fetch(t *testing.T) {
	e := echo.New()
	v := validator.New()

	var articles []ka.Article
	testdata.GoldenJSONUnmarshal(t, "articles", &articles)

	tests := map[string]struct {
		path               string
		fetchArticle       testdata.FuncCaller
		expectedResp       []ka.Article
		expectedStatusCode int
	}{
		"success": {
			path: "/articles",
			fetchArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, ka.ArticleFilter{}},
				Output:   []interface{}{articles, nil},
			},
			expectedResp:       articles,
			expectedStatusCode: http.StatusOK,
		},
		"success with query param": {
			path: "/articles?author=jhon&search=wick",
			fetchArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, ka.ArticleFilter{Author: "jhon", Search: "wick"}},
				Output:   []interface{}{articles, nil},
			},
			expectedResp:       articles,
			expectedStatusCode: http.StatusOK,
		},
		"internal server error": {
			path: "/articles",
			fetchArticle: testdata.FuncCaller{
				IsCalled: true,
				Input:    []interface{}{mock.Anything, ka.ArticleFilter{}},
				Output:   []interface{}{nil, errors.New("error")},
			},
			expectedResp:       nil,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			articleServiceMock := new(mocks.ArticleService)

			if test.fetchArticle.IsCalled {
				articleServiceMock.On("Fetch", test.fetchArticle.Input...).
					Return(test.fetchArticle.Output...).
					Once()
			}

			req := httptest.NewRequest(http.MethodGet, test.path, nil)
			rec := httptest.NewRecorder()

			rest.RegisterPathArticle(e, v, articleServiceMock)
			e.ServeHTTP(rec, req)

			require.Equal(t, test.expectedStatusCode, rec.Code)

			if rec.Code == http.StatusOK {
				var _articles []ka.Article
				err := json.Unmarshal(rec.Body.Bytes(), &_articles)
				require.NoError(t, err)

				require.Equal(t, test.expectedResp, _articles)
			}

			articleServiceMock.AssertExpectations(t)
		})
	}
}
