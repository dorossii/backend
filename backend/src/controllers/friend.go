package controllers

import (
	"backend/services"
	"errors"
	"fmt"
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

func AcceptFriendRequest(ctx echo.Context) error {
	userID := ctx.Get("UserID").(string)

	// 相手のIDを取得
	friendID := ctx.Request().Header.Get("FriendId")
	if friendID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "FriendId header is required"})
	}

	err := services.AcceptFriendRequest(userID, friendID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "success"})
}

// 送られてきたフレンドリクエストを取得する関数
func GetReceivedFriendRequests(ctx echo.Context) error {
	userID := ctx.Get("UserID").(string)

	// 相手のIDを取得
	res, err := services.GetFriendRequests(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, res)
}


func GetFriends(ctx echo.Context) error {
	userID := ctx.Get("UserID").(string)

	friends, err := services.GetFriends(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"friends": friends})
}

func GetInviteURL(ctx echo.Context) error {
	userID := ctx.Get("UserID").(string)

	url := fmt.Sprintf("%s/app/friend/invite?inviteid=%s", inviteBaseURL, userID)

	return ctx.JSON(http.StatusOK, map[string]string{
		"URL": url,
	})
}
