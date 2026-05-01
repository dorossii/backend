package main

// エンドポイントのルーティング
import (
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
	user := router.Group("/user")
	{
		user.POST("/register", TempController)
	}


	return router
}