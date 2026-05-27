package services

import (
	"backend/logger"
	"backend/models"
	"backend/repositories"
	"errors"
)

var (
	ErrAlreadySent     = errors.New("すでにフレンド申請済みです")
	ErrAlreadyReceived = errors.New("相手からすでに申請が届いています")
	ErrFriendNotFound = errors.New("指定されたフレンドが存在しません")
)

func SendFriendRequest(userID, friendID string) error {
	existing, err := repositories.GetFriendShipAny(userID, friendID)
	if err != nil {
		return err
	}
	if existing != nil {
		if existing.UserID == userID {
			return ErrAlreadySent
		}
		return ErrAlreadyReceived
	}

	return repositories.CreateFriendShip(&models.FriendShips{
		UserID:   userID,
		FriendID: friendID,
		Status:   models.FriendStatusPending,
	})
}

func PostAttackerSettings(userID string, targetUser string) error {
	// 空文字ならランダム設定
	if targetUser == "" {
		err := repositories.UpdateAttackerSettings(userID, "")

		if err != nil {
			return err
		}
		return nil
	}

	// 指定ユーザーの場合は friend check
	isFriend, err := repositories.IsFriend(userID, targetUser)
	if err != nil {
		return err
	}

	if !isFriend {
		return ErrFriendNotFound
	}

	err = repositories.UpdateAttackerSettings(userID, targetUser)
	if err != nil {
		logger.PrintErr("update attacker settings", err)
		return err
	}

	return nil
}