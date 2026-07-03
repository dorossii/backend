package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RegisterUserRequest struct {
	BirthDate  int64  `json:"birthDate"`
	LivingType string `json:"livingType"`
}

type RegisterUserResponse struct {
	UserID     string `json:"userId"`
	UserName   string `json:"userName"`
	BirthDate  int64  `json:"birthDate"`
	LivingType string `json:"livingType"`
}

func RegisterUser(userID, userName, email string, req RegisterUserRequest) (*RegisterUserResponse, error) {
	if req.LivingType != "alone" && req.LivingType != "family" {
		return nil, errors.New("livingType は alone か family のみ有効です")
	}

	birthDate := time.Unix(req.BirthDate, 0).UTC()

	user := &models.User{
		UserID:     userID,
		UserName:   userName,
		BirthDate:  birthDate,
		Mailadress: email,
	}

	userRoom := &models.UserRoom{
		UserID:  userID,
		IsAlone: req.LivingType == "alone",
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if err := repositories.CreateUserTx(tx, user); err != nil {
			return err
		}
		return repositories.CreateUserRoomTx(tx, userRoom)
	})
	if err != nil {
		return nil, err
	}

	return &RegisterUserResponse{
		UserID:     user.UserID,
		UserName:   user.UserName,
		BirthDate:  user.BirthDate.Unix(),
		LivingType: req.LivingType,
	}, nil
}

type LifestyleRequest struct {
	IsAlone   bool `json:"isAlone"`
	HasWasher bool `json:"hasWasher"`
	HasVacuum bool `json:"hasVacuum"`
}

type LifestyleResponse struct {
	IsAlone   bool `json:"isAlone"`
	HasWasher bool `json:"hasWasher"`
	HasVacuum bool `json:"hasVacuum"`
}

// CreateUserLifestyle は生活環境情報を登録する（初回）
// register時点で UserRoom が既に作成されている可能性があるため UPSERT する
func CreateUserLifestyle(userID string, req LifestyleRequest) (*LifestyleResponse, error) {
	if err := repositories.UpsertUserRoom(userID, req.IsAlone, req.HasWasher, req.HasVacuum); err != nil {
		return nil, err
	}

	return &LifestyleResponse{
		IsAlone:   req.IsAlone,
		HasWasher: req.HasWasher,
		HasVacuum: req.HasVacuum,
	}, nil
}

// UpdateUserLifestyle は生活環境情報を編集する（2回目以降）
func UpdateUserLifestyle(userID string, req LifestyleRequest) error {
	return repositories.UpdateUserRoom(userID, req.IsAlone, req.HasWasher, req.HasVacuum)
}
