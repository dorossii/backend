package controllers

import (
	"backend/models"
	"backend/services"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminListTasks(ctx echo.Context) error {
	userID := ctx.QueryParam("userID")

	tasks, err := services.AdminListTasks(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"tasks": tasks})
}

func AdminCreateTask(ctx echo.Context) error {
	var task models.Task
	if err := ctx.Bind(&task); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := services.AdminCreateTask(&task); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, task)
}

func AdminUpdateTask(ctx echo.Context) error {
	taskID := ctx.Param("id")

	var task models.Task
	if err := ctx.Bind(&task); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	updated, err := services.AdminUpdateTask(taskID, &task)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, updated)
}

type AdminUpdateTaskStatusRequest struct {
	Status models.TaskStatus `json:"Status"`
}

func AdminUpdateTaskStatus(ctx echo.Context) error {
	taskID := ctx.Param("id")

	var req AdminUpdateTaskStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	updated, err := services.AdminUpdateTaskStatus(taskID, req.Status)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, updated)
}

func AdminDeleteTask(ctx echo.Context) error {
	taskID := ctx.Param("id")

	if err := services.AdminDeleteTask(taskID); err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func AdminGetTaskImage(ctx echo.Context) error {
	taskID := ctx.Param("id")

	filePath, err := services.AdminGetTaskImagePath(taskID)
	if err != nil {
		if errors.Is(err, services.ErrTaskNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.File(filePath)
}
