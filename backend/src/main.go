package main

import (
	"backend/batch"
	"backend/middlewares"
	"backend/models"
	"backend/services"
	"backend/seeds"
	"log"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// モデル初期化
	models.Init()

	// シードデータ投入
	err := seeds.Seed()
	if err != nil {
		log.Fatal("Failed to seed data: ", err)
	}

	// task画像サービス初期化
	if err := services.InitTaskImageService(); err != nil {
		log.Fatal(err)
	}

	// バッチ処理開始
	batch.Run()
	
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	// ミドルウェア初期化
	middlewares.Init()

	// 管理画面用セッションストアの初期化
	middlewares.InitSession()

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
