package services

import (
	"backend/models"
	"backend/repositories"
	"errors"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("指定されたユーザーが存在しません")

func AdminListUsers(search string) ([]models.User, error) {
	return repositories.ListUsers(search)
}

func AdminGetUser(userID string) (*models.User, error) {
	user, err := repositories.GetUser(userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return user, err
}

func AdminUpdateUserName(userID string, userName string) (*models.User, error) {
	if userName == "" {
		return nil, errors.New("userName は必須です")
	}

	if _, err := AdminGetUser(userID); err != nil {
		return nil, err
	}

	if err := repositories.UpdateUserName(userID, userName); err != nil {
		return nil, err
	}

	return AdminGetUser(userID)
}

type AdminUserStatsUpdate struct {
	HealthPoint *int `json:"HealthPoint"`
	DirtLevel   *int `json:"DirtLevel"`
}

func AdminUpdateUserStats(userID string, update AdminUserStatsUpdate) (*models.User, error) {
	if _, err := AdminGetUser(userID); err != nil {
		return nil, err
	}

	if err := repositories.UpdateUserStats(userID, update.HealthPoint, update.DirtLevel); err != nil {
		return nil, err
	}

	return AdminGetUser(userID)
}
