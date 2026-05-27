package repositories

import (
	"backend/models"
	"errors"

	"gorm.io/gorm"
)

func IsFriend(userId string, friendId string) (bool, error) {
	var friend models.FriendShips

	err := models.DB.
		Where(
			"user_id = ? AND friend_id = ? AND status = ?",
			userId,
			friendId,
			1, // フレンド関係が成立している
		).
		First(&friend).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err // エラーが出てるのにfalse返すのきしょいかも
	}

	return true, nil
}
