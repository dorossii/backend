package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
	"time"
)

type RegisterUserRequest struct {
	UserName    string `json:"userName"`
	BirthDate   int64  `json:"birthDate"`
	MailAddress string `json:"mailAddress"`
	LivingType  string `json:"livingType"`
}

type RegisterUserResponse struct {
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	BirthDate   int64  `json:"birthDate"`
	MailAddress string `json:"mailAddress"`
	LivingType  string `json:"livingType"`
}

func RegisterUser(userID string, req RegisterUserRequest) (*RegisterUserResponse, error) {
	if req.LivingType != "alone" && req.LivingType != "together" {
		return nil, errors.New("livingType は alone か together のみ有効です")
	}

	birthDate := time.Unix(req.BirthDate, 0).UTC()

	user := &models.User{
		UserID:     userID,
		UserName:   req.UserName,
		BirthDate:  birthDate,
		Mailadress: req.MailAddress,
	}

	if err := repositories.CreateUser(user); err != nil {
		return nil, err
	}

	return &RegisterUserResponse{
		UserID:      user.UserID,
		UserName:    user.UserName,
		BirthDate:   user.BirthDate.Unix(),
		MailAddress: user.Mailadress,
		LivingType:  req.LivingType,
	}, nil
}
