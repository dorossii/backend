package controllers

import (
	"backend/models"
	"backend/services"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminListBaseTasks(ctx echo.Context) error {
	baseTasks, err := services.AdminListBaseTasks()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"baseTasks": baseTasks})
}

func AdminCreateBaseTask(ctx echo.Context) error {
	var bt models.BaseTask
	if err := ctx.Bind(&bt); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := services.AdminCreateBaseTask(&bt); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, bt)
}

func AdminUpdateBaseTask(ctx echo.Context) error {
	baseID := ctx.Param("id")

	var bt models.BaseTask
	if err := ctx.Bind(&bt); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	updated, err := services.AdminUpdateBaseTask(baseID, &bt)
	if err != nil {
		if errors.Is(err, services.ErrBaseTaskNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, updated)
}

func AdminDeleteBaseTask(ctx echo.Context) error {
	baseID := ctx.Param("id")

	if err := services.AdminDeleteBaseTask(baseID); err != nil {
		if errors.Is(err, services.ErrBaseTaskNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
