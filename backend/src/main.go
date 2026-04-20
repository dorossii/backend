package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	router := echo.New()
	router.Use(middleware.RequestLogger())
	router.Use(middleware.Recover())

	router.GET("/", func(ctx *echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, World!")
	})

	if err := router.Start("0.0.0.0:8090"); err != nil {
		router.Logger.Error("failed to start server", "error", err)
	}
}
