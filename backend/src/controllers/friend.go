package controllers

import (
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
