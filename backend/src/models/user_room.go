package models

type UserRoom struct {
	UserID    string `json:"UserID" gorm:"primaryKey"`      // ユーザーID
	IsAlone   bool   `json:"IsAlone" gorm:"default:false"`  // 一人暮らしか
	HasWasher bool   `json:"HasWasher" gorm:"default:true"` // 洗濯機があるか
	HasVacuum bool   `json:"HasVacuum" gorm:"default:true"` // 掃除機があるか
}