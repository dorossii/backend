package controllers

import (
	"backend/services"
	"net/http"

	"github.com/labstack/echo/v4"
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
