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

// GetFriends は userID の承認済みフレンドの User 情報一覧を返す
func GetFriends(userID string) ([]*models.User, error) {
	var users []*models.User
	err := models.DB.
		Joins("JOIN friend_ships ON (friend_ships.friend_id = users.user_id OR friend_ships.user_id = users.user_id)").
		Where("(friend_ships.user_id = ? OR friend_ships.friend_id = ?) AND friend_ships.status = ? AND users.user_id != ?",
			userID, userID, models.FriendStatusAccepted, userID).
		Find(&users).Error
	return users, err
}
