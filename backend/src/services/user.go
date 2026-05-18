package services

import (
	"backend/models"
	"backend/repositories"
	"errors"
	"time"
)

type RegisterUserRequest struct {
	UserName    string `json:"userName"`
	BirthDate   string `json:"birthDate"`
	MailAddress string `json:"mailAddress"`
	LivingType  string `json:"livingType"`
}

type RegisterUserResponse struct {
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	BirthDate   string `json:"birthDate"`
	MailAddress string `json:"mailAddress"`
	LivingType  string `json:"livingType"`
}

func RegisterUser(userID string, req RegisterUserRequest) (*RegisterUserResponse, error) {
	if req.LivingType != "alone" && req.LivingType != "together" {
		return nil, errors.New("livingType は alone か together のみ有効です")
	}

	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return nil, errors.New("birthDate の形式が不正です (YYYY-MM-DD)")
	}

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
		BirthDate:   user.BirthDate.Format("2006-01-02"),
		MailAddress: user.Mailadress,
		LivingType:  req.LivingType,
	}, nil
}
