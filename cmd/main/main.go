package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	
	"github.com/nugrohoac/kumparan-assessment/cmd"
	"github.com/nugrohoac/kumparan-assessment/internal/rest"
)

func main() {
	e := echo.New()
	v := validator.New()

	e.Use(echoMiddleware.CORS())

	rest.RegisterPathArticle(e, v, cmd.ArticleService)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cmd.PortApp)))
}
