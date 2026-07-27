package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// AdminRateLimiter は管理画面(ブルートフォース対策)向けに、IPアドレスごとに1秒間5リクエストへ制限するミドルウェア
var AdminRateLimiter = middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{
		Rate:  rate.Limit(5), // 1秒あたりの許容リクエスト数
		Burst: 10,
	}),
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		return ctx.RealIP(), nil
	},
})
