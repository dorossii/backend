package controllers

import (
	"backend/middlewares"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type AdminLoginRequest struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func AdminLogin(ctx echo.Context) error {
	var req AdminLoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if req.Username == "" || req.Password == "" ||
		req.Username != os.Getenv("ADMIN_USERNAME") || req.Password != os.Getenv("ADMIN_PASSWORD") {
		return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	session, err := middlewares.Store.Get(ctx.Request(), middlewares.AdminSessionName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	session.Values["authenticated"] = true
	if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func AdminLogout(ctx echo.Context) error {
	session, err := middlewares.Store.Get(ctx.Request(), middlewares.AdminSessionName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	if err := session.Save(ctx.Request(), ctx.Response()); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
