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

func UpdateFriendShip(fs *models.FriendShips) error {
	return models.DB.Save(fs).Error
}

func GetFriendShipsByStatus(userID string, status models.FriendStatus) ([]*models.FriendShips, error) {
	var fs []*models.FriendShips
	err := models.DB.Where("user_id = ? AND status = ?", userID, status).Find(&fs).Error
	return fs, err
}

func GetIncomingFriendShipsByStatus(userID string, status models.FriendStatus) ([]*models.FriendShips, error) {
	var fs []*models.FriendShips
	err := models.DB.Where("friend_id = ? AND status = ?", userID, status).Find(&fs).Error
	return fs, err
}

func GetRescueUserIDs(userID string)([]models.HelpTargets, error) {
var helpTargets []models.HelpTargets

	err := models.DB.Where("user_id = ?", userID).Find(&helpTargets).Error
	if err != nil {
		return nil, err
	}

	return helpTargets, nil
}