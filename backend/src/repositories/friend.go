package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

func GetFriendShip(userID, friendID string) (*models.FriendShips, error) {
	var fs models.FriendShips
	err := models.DB.First(&fs, "user_id = ? AND friend_id = ?", userID, friendID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &fs, err
}

// GetFriendShipAny は userID と friendID の双方向いずれかのレコードを返す
func GetFriendShipAny(userID, friendID string) (*models.FriendShips, error) {
	if fs, err := GetFriendShip(userID, friendID); err != nil || fs != nil {
		return fs, err
	}
	return GetFriendShip(friendID, userID)
}

func CreateFriendShip(fs *models.FriendShips) error {
	return models.DB.Create(fs).Error
}
