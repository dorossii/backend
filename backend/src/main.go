package main

import (
	"backend/middlewares"
	"backend/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// モデル初期化
	models.Init()
	
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// ミドルウェア初期化
	middlewares.Init()

	// ルーティングの設定を追加
	router = InitRouter(router)

	router.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "Hello, World!")
	})

	router.GET("/authed", middlewares.RequireAuth(func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK,map[string]string{"message": "authed"})
	}))

	if err := router.Start("0.0.0.0:8090"); err != nil {
		router.Logger.Error("failed to start server", "error", err)
	}
}
