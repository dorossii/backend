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

	if existing, err := repositories.GetUser(userID); err == nil {
		// 既に登録済みの場合は何もせず成功として返す（認証機からのprovisioning呼び出しが複数回発生しうるため）
		return &RegisterUserResponse{
			UserID:     existing.UserID,
			UserName:   existing.UserName,
			BirthDate:  existing.BirthDate.Unix(),
			LivingType: req.LivingType,
		}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	birthDate := time.Unix(req.BirthDate, 0).UTC()

	user := &models.User{
		UserID:     userID,
		UserName:   userName,
		BirthDate:  birthDate,
		Mailadress: email,
		Icon:       "pineTree",
		BgColor:    "icon1",
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

	// 最初のタスクをユーザーに割り当てる
	if err := repositories.CreateTaskForUser(userID); err != nil {
		return nil, err
	}

	return &RegisterUserResponse{
		UserID:     user.UserID,
		UserName:   user.UserName,
		BirthDate:  user.BirthDate.Unix(),
		LivingType: req.LivingType,
	}, nil
}

// UserStatusResponse はホーム画面向けのユーザーステータス
type UserStatusResponse struct {
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	DirtLevel   int    `json:"DirtLevel"`
	HealthPoint int    `json:"HealthPoint"`
}

// GetUserStatus はホーム画面向けのユーザーステータスを取得する
func GetUserStatus(userID string) (*UserStatusResponse, error) {
	user, err := repositories.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return &UserStatusResponse{
		UserID:      user.UserID,
		UserName:    user.UserName,
		DirtLevel:   user.DirtLevel,
		HealthPoint: user.HealthPoint,
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

// UserSettingRequest は PUT /app/user/setting のリクエストボディ
type UserSettingRequest struct {
	UserName string  `json:"userName"`         // 変更後のユーザー名
	Icon     *string `json:"icon"`             // 変更後のアイコン識別子(任意)
	BgColor  *string `json:"background"`       // 変更後の背景識別子(任意)
}

// UpdateUserSetting はユーザー名・アイコン・背景を更新する
// Icon/BgColor はリクエストに含まれない場合(nil)は既存値を維持する
func UpdateUserSetting(userID string, req UserSettingRequest) error {
	if req.UserName == "" { // 空文字は許可しない
		return errors.New("userName は必須です")
	}

	return repositories.UpdateUserSetting(userID, req.UserName, req.Icon, req.BgColor) // リポジトリ経由でDBを更新する
}
