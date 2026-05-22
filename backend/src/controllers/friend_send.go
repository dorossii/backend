package controllers

import (
	"backend/services"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendFriendRequest(ctx echo.Context) error {
	userID := ctx.Get("UserID").(string)

	friendID := ctx.Request().Header.Get("InviteID")
	if friendID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "InviteID header is required"})
	}

	err := services.SendFriendRequest(userID, friendID)
	if err != nil {
		if errors.Is(err, services.ErrAlreadySent) || errors.Is(err, services.ErrAlreadyReceived) {
			return ctx.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
