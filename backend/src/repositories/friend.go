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

func CreateFriendShip(fs *models.FriendShips) error {
	return models.DB.Create(fs).Error
}
