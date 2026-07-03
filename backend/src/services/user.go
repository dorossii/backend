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
		UserID:        userID,                    // ユーザーID
		IsAlone:       req.LivingType == "alone", // livingTypeから一人暮らしかどうかを判定する
		HasWasher:     true,                      // register時点ではデフォルト値で作成する
		HasVacuum:     true,                      // register時点ではデフォルト値で作成する
		HasRobot:      true,                      // register時点ではデフォルト値で作成する
		UseTableware:  true,                      // register時点ではデフォルト値で作成する
		HasDishwasher: true,                      // register時点ではデフォルト値で作成する
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

// LifestyleRequest は生活環境情報の登録・編集リクエスト
type LifestyleRequest struct {
	IsAlone       bool `json:"isAlone"`       // 一人暮らしか
	HasWasher     bool `json:"hasWasher"`     // 洗濯機があるか
	HasVacuum     bool `json:"hasVacuum"`     // 掃除機があるか
	HasRobot      bool `json:"hasRobot"`      // ロボット掃除機があるか
	UseTableware  bool `json:"useTableware"`  // 食器を使用するか
	HasDishwasher bool `json:"hasDishwasher"` // 食洗機があるか
}

// LifestyleResponse は生活環境情報のレスポンス
type LifestyleResponse struct {
	IsAlone       bool `json:"isAlone"`       // 一人暮らしか
	HasWasher     bool `json:"hasWasher"`     // 洗濯機があるか
	HasVacuum     bool `json:"hasVacuum"`     // 掃除機があるか
	HasRobot      bool `json:"hasRobot"`      // ロボット掃除機があるか
	UseTableware  bool `json:"useTableware"`  // 食器を使用するか
	HasDishwasher bool `json:"hasDishwasher"` // 食洗機があるか
}

// lifestyleRequestToUserRoom は LifestyleRequest を models.UserRoom に変換する
func lifestyleRequestToUserRoom(userID string, req LifestyleRequest) models.UserRoom {
	return models.UserRoom{
		UserID:        userID,            // ユーザーID
		IsAlone:       req.IsAlone,       // 一人暮らしか
		HasWasher:     req.HasWasher,     // 洗濯機があるか
		HasVacuum:     req.HasVacuum,     // 掃除機があるか
		HasRobot:      req.HasRobot,      // ロボット掃除機があるか
		UseTableware:  req.UseTableware,  // 食器を使用するか
		HasDishwasher: req.HasDishwasher, // 食洗機があるか
	}
}

// CreateUserLifestyle は生活環境情報を登録する（初回）
// register時点で UserRoom が既に作成されている可能性があるため UPSERT する
func CreateUserLifestyle(userID string, req LifestyleRequest) (*LifestyleResponse, error) {
	room := lifestyleRequestToUserRoom(userID, req) // リクエストをUserRoomに変換する

	if err := repositories.UpsertUserRoom(userID, room); err != nil {
		return nil, err
	}

	return &LifestyleResponse{
		IsAlone:       req.IsAlone,       // 一人暮らしか
		HasWasher:     req.HasWasher,     // 洗濯機があるか
		HasVacuum:     req.HasVacuum,     // 掃除機があるか
		HasRobot:      req.HasRobot,      // ロボット掃除機があるか
		UseTableware:  req.UseTableware,  // 食器を使用するか
		HasDishwasher: req.HasDishwasher, // 食洗機があるか
	}, nil
}

// UpdateUserLifestyle は生活環境情報を編集する（2回目以降）
func UpdateUserLifestyle(userID string, req LifestyleRequest) error {
	room := lifestyleRequestToUserRoom(userID, req) // リクエストをUserRoomに変換する

	return repositories.UpdateUserRoom(userID, room)
}

type UserSettingRequest struct {
	UserName string `json:"userName"`
}

func UpdateUserSetting(userID string, req UserSettingRequest) error {
	if req.UserName == "" {
		return errors.New("userName は必須です")
	}

	return repositories.UpdateUserName(userID, req.UserName)
}
