package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}

func CreateUserTx(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

// CreateUserRoomTx は UserRoom を作成する。
// bool フィールドに gorm:"default:true" が付与されているため、構造体経由の Create だと
// ゼロ値(false)が「未設定」とみなされデフォルト値で上書きされてしまう。
// map 経由で書き込むことでゼロ値も確実に反映する。
func CreateUserRoomTx(tx *gorm.DB, userRoom *models.UserRoom) error {
	return tx.
		Model(&models.UserRoom{}).
		Create(map[string]interface{}{
			"user_id":    userRoom.UserID,
			"is_alone":   userRoom.IsAlone,
			"has_washer": userRoom.HasWasher,
			"has_vacuum": userRoom.HasVacuum,
		}).Error
}

// UpsertUserRoom は UserRoom を UPSERT する（存在すれば全カラム更新、なければ作成）
func UpsertUserRoom(userID string, isAlone, hasWasher, hasVacuum bool) error {
	var count int64
	if err := models.DB.Model(&models.UserRoom{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return models.DB.
			Model(&models.UserRoom{}).
			Create(map[string]interface{}{
				"user_id":    userID,
				"is_alone":   isAlone,
				"has_washer": hasWasher,
				"has_vacuum": hasVacuum,
			}).Error
	}

	return UpdateUserRoom(userID, isAlone, hasWasher, hasVacuum)
}

func UpdateUserRoom(userID string, isAlone, hasWasher, hasVacuum bool) error {
	return models.DB.
		Model(&models.UserRoom{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_alone":   isAlone,
			"has_washer": hasWasher,
			"has_vacuum": hasVacuum,
		}).Error
}

func GetUser(UserID string) (*models.User, error) {
	var user models.User
	err := models.DB.First(&user, "user_id = ?", UserID).Error
	return &user, err
}

func UpdateAttackerSettings(userID string, targetUser string) error {
	err := models.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("target_user", targetUser).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateAttackerSettingsTx(tx *gorm.DB, userID string, targetUser string) error {
	err := tx.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("target_user", targetUser).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAttackerSettings は userID の target_user を返す	
func GetAttackerSettings(userID string) (string, error) {
	var user models.User
	err := models.DB.Select("target_user").First(&user, "user_id = ?", userID).Error
	if err != nil {
		return "", err
	}
	
	return user.TargetUser, nil
}

func UpdateDirtLevel(tx *gorm.DB, userID string, diff int) error {
	var user models.User

	err := tx.Model(&models.User{}).First(&user, "user_id = ?", userID).Error
	if err != nil {
		return err
	}

	// 計算後の値が0以下700以上ならぱわーでねじ曲げる
	newDirt := user.DirtLevel + diff

	if newDirt < 0 {
		newDirt = 0
	}

	if newDirt > 700 {
		newDirt = 700
	}

	return tx.Model(&models.User{}).Where("user_id = ?", userID).Update("dirt_level", newDirt).Error
}

