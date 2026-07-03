package models

// UserRoom はユーザーの生活環境情報を表す
type UserRoom struct {
	UserID        string `json:"UserID" gorm:"primaryKey"`          // ユーザーID
	IsAlone       bool   `json:"IsAlone" gorm:"default:false"`      // 一人暮らしか（false:実家暮らし）
	HasWasher     bool   `json:"HasWasher" gorm:"default:true"`     // 洗濯機があるか
	HasVacuum     bool   `json:"HasVacuum" gorm:"default:true"`     // 掃除機があるか
	HasRobot      bool   `json:"HasRobot" gorm:"default:true"`      // ロボット掃除機があるか
	UseTableware  bool   `json:"UseTableware" gorm:"default:true"`  // 食器を使用するか
	HasDishwasher bool   `json:"HasDishwasher" gorm:"default:true"` // 食洗機があるか
}
