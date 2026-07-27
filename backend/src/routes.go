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

		// ホーム画面 - ユーザーステータス取得
		user.GET("/status", controllers.GetUserStatus)

		// ユーザー情報の編集
		user.PUT("/setting", controllers.UpdateUserSetting)

		// 生活環境情報の登録
		user.POST("/lifestyle", controllers.CreateUserLifestyle)

		// 生活環境情報の編集
		user.PUT("/lifestyle", controllers.UpdateUserLifestyle)

		// タスクグループ
		task := user.Group("/task")
		{
			// タスク取得
			task.GET("", controllers.GetTask)

			// タスク写真のアップロード
			task.POST("/:id/image", controllers.PostUploadImageHandler)

			// タスク煽りメッセージ
			task.POST("/message", controllers.PostTauntMessageHandler)

			// タスクのステータス更新
			task.PUT("/:id", controllers.PutTaskStatusHandler)
		}

		// 承認待ちタスク一覧取得
		user.GET("/tasks/pending", controllers.GetPendingTasksHandler)

		// タスク複数完了
		user.POST("/tasks/complete", TempController)
	}

	// タスク写真確認（認証不要）
	router.GET("/user/task/:imageId/image", controllers.GetTaskImageHandler)

	// friendグループ
	friend := router.Group("/friend", middlewares.RequireAuth)
	{
		// フレンド一覧取得
		friend.GET("", controllers.GetFriends)

		// フレンド招待
		friend.GET("/invite", controllers.GetInviteURL)

		// フレンド申請送信
		friend.POST("/send", controllers.SendFriendRequest)

		// フレンド認証
		friend.POST("/accept", controllers.AcceptFriendRequest)

		// フレンドリクエスト一覧を取得
		friend.GET("/received", controllers.GetReceivedFriendRequests)

		// フレンド削除
		friend.DELETE("/:id", controllers.DeleteFriend)

		// 嫌がらせする人の取得
		friend.GET("/attack", controllers.GetAttackerSettingsHandler)

		// 嫌がらせする人の設定
		friend.PUT("/attack", controllers.PostAttackerSettingsHandler)

		// レスキューする人の設定
		friend.POST("/rescue", controllers.PostRescuerSettingsHandler)

		// レスキューしてほしい人の一覧取得（isrescued: help_targets に登録済みなら true）
		friend.GET("/rescue", controllers.GetRescueFriends)
	}

	// noticeグループ
	notice := router.Group("/notice")
	{
		// 通知取得
		notice.GET("/", TempController)
	}

	// 管理画面ログイン・ログアウト(認証不要、ブルートフォース対策でレート制限)
	router.POST("/admin/login", controllers.AdminLogin, middlewares.AdminRateLimiter)
	router.POST("/admin/logout", controllers.AdminLogout, middlewares.AdminRateLimiter)

	// 管理画面 静的アセット(CSS/JS, 認証不要)
	router.Static("/admin/static", "admin_ui")

	// 管理画面ログインページ(認証不要)
	router.GET("/admin/login.html", func(ctx echo.Context) error {
		return ctx.File("admin_ui/login.html")
	})

	// 管理画面 本体(セッション認証必須、未認証時はログインページへリダイレクト)
	router.GET("/admin/", func(ctx echo.Context) error {
		return ctx.File("admin_ui/index.html")
	}, middlewares.RequireAdminSessionPage, middlewares.AdminRateLimiter)

	// adminグループ(ブルートフォース対策でIPごとに1秒1リクエストへ制限)
	admin := router.Group("/admin", middlewares.RequireAdminSession, middlewares.AdminRateLimiter)
	{
		// BaseTask(タスクテンプレート) CRUD
		baseTask := admin.Group("/base-tasks")
		{
			baseTask.GET("", controllers.AdminListBaseTasks)
			baseTask.POST("", controllers.AdminCreateBaseTask)
			baseTask.PUT("/:id", controllers.AdminUpdateBaseTask)
			baseTask.DELETE("/:id", controllers.AdminDeleteBaseTask)
		}

		// Task(個別タスク)管理
		task := admin.Group("/tasks")
		{
			task.GET("", controllers.AdminListTasks)
			task.POST("", controllers.AdminCreateTask)
			task.PUT("/:id", controllers.AdminUpdateTask)
			task.PUT("/:id/status", controllers.AdminUpdateTaskStatus)
			task.DELETE("/:id", controllers.AdminDeleteTask)
			task.GET("/:id/image", controllers.AdminGetTaskImage)
		}

		// フレンド関係管理
		friendAdmin := admin.Group("/friendships")
		{
			friendAdmin.POST("", controllers.AdminCreateFriendShip)
			friendAdmin.PUT("/:userId/:friendId", controllers.AdminUpdateFriendShipStatus)
			friendAdmin.DELETE("/:userId/:friendId", controllers.AdminDeleteFriendShip)
		}

		// ユーザー管理(HP/DirtLevel編集含む)
		userAdmin := admin.Group("/users")
		{
			userAdmin.GET("", controllers.AdminListUsers)
			userAdmin.GET("/:id", controllers.AdminGetUser)
			userAdmin.PUT("/:id/stats", controllers.AdminUpdateUserStats)
		}
	}

	return router
}
