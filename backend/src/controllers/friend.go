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


// GetFriends は承認済みフレンドの一覧を返すハンドラ
func GetFriends(ctx echo.Context) error {
	// JWTミドルウェアで検証済みのユーザーIDを取得
	userID := ctx.Get("UserID").(string)

	// フレンド一覧を取得
	friends, err := services.GetFriends(userID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"friends": friends})
}

// DeleteFriend はフレンド関係を削除するハンドラ
func DeleteFriend(ctx echo.Context) error {
	// JWTミドルウェアで検証済みのユーザーIDを取得
	userID := ctx.Get("UserID").(string)

	// パスパラメータから削除対象のフレンドIDを取得
	friendID := ctx.Param("id")
	if friendID == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	// フレンド関係を削除
	err := services.DeleteFriend(userID, friendID)
	if err != nil {
		if errors.Is(err, services.ErrFriendShipNotFound) || errors.Is(err, services.ErrFriendShipNotAccepted) {
			return ctx.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}

func GetInviteURL(ctx echo.Context) error {
	userID := ctx.Get("UserID").(string)

	url := fmt.Sprintf("%s/app/friend/invite?inviteid=%s", inviteBaseURL, userID)

	return ctx.JSON(http.StatusOK, map[string]string{
		"URL": url,
	})
}
