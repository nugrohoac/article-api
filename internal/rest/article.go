package rest

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	ka "github.com/nugrohoac/kumparan-assessment"
)

type articleHandler struct {
	service ka.ArticleService
	v       *validator.Validate
}

func (a articleHandler) store(c echo.Context) error {
	var article ka.Article
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := a.v.Struct(article); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	article, err := a.service.Store(c.Request().Context(), article)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, 0)
	}

	return c.JSON(http.StatusCreated, article)
}

func (a articleHandler) fetch(c echo.Context) error {
	filter := ka.ArticleFilter{
		Author: c.QueryParam("author"),
		Search: c.QueryParam("search"),
	}

	articles, err := a.service.Fetch(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, articles)
}

// RegisterPathArticle .
func RegisterPathArticle(e *echo.Echo, v *validator.Validate, service ka.ArticleService) {
	handler := articleHandler{
		service: service,
		v:       v,
	}

	e.GET("/articles", handler.fetch)
	e.POST("/articles", handler.store)
}
