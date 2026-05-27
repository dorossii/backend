package controllers

import (
	"backend/logger"
	"backend/services"

	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetInviteURL(ctx echo.Context) error {
	userID := ctx.Get("UserID").(string)

	url := fmt.Sprintf("%s/app/friend/invite?inviteid=%s", inviteBaseURL, userID)

	return ctx.JSON(http.StatusOK, map[string]string{
		"URL": url,
	})
}

type PostAttackerSettingsRequest struct {
	TargetUser string `json:"target_user"`
}

// 嫌がらせする人の設定
func PostAttackerSettingsHandler(ctx echo.Context) error {
	id := ctx.Request().Header.Get("UserID")

	var req PostAttackerSettingsRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid request",
		})
	}

	err := services.PostAttackerSettings(id, req.TargetUser)
	if err != nil {
		logger.PrintErr("PostAttackerSettingsHandler", err)

		if errors.Is(err, services.ErrFriendNotFound) {
			return ctx.JSON(http.StatusForbidden, echo.Map{
				"error": "friend not found",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})

}
