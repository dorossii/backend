package controllers

import (
	"backend/services"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterUser(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)
	name := ctx.Get("Name").(string)
	email := ctx.Get("Email").(string)

	var req services.RegisterUserRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	res, err := services.RegisterUser(userId, name, email, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}

func GetUserStatus(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	res, err := services.GetUserStatus(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}

func CreateUserLifestyle(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	var req services.LifestyleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	res, err := services.CreateUserLifestyle(userId, req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}

func UpdateUserLifestyle(ctx echo.Context) error {
	userId := ctx.Get("UserID").(string)

	var req services.LifestyleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := services.UpdateUserLifestyle(userId, req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{})
}
