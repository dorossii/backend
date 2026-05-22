package main

// エンドポイントのルーティング
import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

// エラー避け
func TempController(c echo.Context) error {
	return c.JSON(501, map[string]string{
		"message": "not implemented",
	})
}

// ルーティング
func InitRouter(router *echo.Echo) *echo.Echo {
	// userグループ
	user := router.Group("/user", middlewares.RequireAuth)
	{
		// 初回ユーザー登録
		user.POST("/register", controllers.RegisterUser)

		// ユーザー情報の編集
		user.PUT("/setting", TempController)

		// 生活環境情報の登録
		user.POST("/lifestyle", TempController)

		// 生活環境情報の編集
		user.PUT("/lifestyle", TempController)

		// タスクグループ
		task := user.Group("/task")
		{
			// タスク取得
			task.GET("", TempController)

			// タスク詳細取得
			task.GET("/:id", TempController)

			// タスク写真のアップロード
			task.POST("/:id/image", TempController)

			// タスク煽りメッセージ
			task.POST("/:id/message", TempController)

			// 写真確認
			task.GET("/:id/image", TempController)

			// タスクのステータス更新
			task.PUT("/:task_id", TempController)
		}

		// タスク複数完了
		user.POST("/tasks/complete", TempController)
	}

	// friendグループ
	friend := router.Group("/friend", middlewares.RequireAuth)
	{
		// フレンド一覧取得
		friend.GET("", TempController)

		// フレンド招待
		friend.GET("/invite", controllers.GetInviteURL)

		// フレンド認証
		friend.POST("/accept", TempController)

		// フレンド削除
		friend.DELETE("/:id", TempController)

		// 嫌がらせする人の設定
		friend.PUT("/attack", TempController)

		// レスキューする人の設定
		friend.POST("/rescue", TempController)
	}

	// noticeグループ
	notice := router.Group("/notice")
	{
		// 通知取得
		notice.GET("/", TempController)
	}

	return router
}
