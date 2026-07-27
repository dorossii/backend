package controllers

import (
	"backend/models"
	"backend/services"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AdminFriendShipRequest struct {
	UserID   string `json:"UserID"`
	FriendID string `json:"FriendID"`
}

func AdminCreateFriendShip(ctx echo.Context) error {
	var req AdminFriendShipRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := services.AdminForceFriendShip(req.UserID, req.FriendID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

type AdminUpdateFriendShipStatusRequest struct {
	Status models.FriendStatus `json:"Status"`
}

func AdminUpdateFriendShipStatus(ctx echo.Context) error {
	userID := ctx.Param("userId")
	friendID := ctx.Param("friendId")

	var req AdminUpdateFriendShipStatusRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := services.AdminUpdateFriendShipStatus(userID, friendID, req.Status)
	if err != nil {
		if errors.Is(err, services.ErrFriendShipNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

func AdminDeleteFriendShip(ctx echo.Context) error {
	userID := ctx.Param("userId")
	friendID := ctx.Param("friendId")

	err := services.AdminDeleteFriendShip(userID, friendID)
	if err != nil {
		if errors.Is(err, services.ErrFriendShipNotFound) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}
