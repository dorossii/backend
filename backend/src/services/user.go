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
