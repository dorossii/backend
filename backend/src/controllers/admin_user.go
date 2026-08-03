package controllers

import (
	"backend/services"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminListUsers(ctx echo.Context) error {
	search := ctx.QueryParam("search")

	users, err := services.AdminListUsers(search)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"users": users})
}

func AdminGetUser(ctx echo.Context) error {
	userID := ctx.Param("id")

	user, err := services.AdminGetUser(userID)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func AdminUpdateUserName(ctx echo.Context) error {
	userID := ctx.Param("id")

	var req struct {
		UserName string `json:"UserName"`
	}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	user, err := services.AdminUpdateUserName(userID, req.UserName)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}

func AdminUpdateUserStats(ctx echo.Context) error {
	userID := ctx.Param("id")

	var req services.AdminUserStatsUpdate
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	user, err := services.AdminUpdateUserStats(userID, req)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, user)
}
