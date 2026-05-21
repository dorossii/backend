package batch

import (
	"backend/models"
	"time"

	"gorm.io/gorm"
)


func DirtTicker() {
	go func() {
		//1時間ごとにタスクを作成する
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				AddDirt()
			}
		}
	}()
}

// 全ユーザーのDirtLevelを3増加させる
func AddDirt() error {
	//汚さが700を超えていないユーザーのDirtLevelを3増加させる
	if err := models.DB.Model(&models.User{}).Where("dirt_level < ?", 700).Update("dirt_level", gorm.Expr("dirt_level + ?", 3)).Error; err != nil {
		return err
	}

	//汚さが700を超えたユーザーはDirtLevelを700に固定する
	if err := models.DB.Model(&models.User{}).Where("dirt_level >= ?", 700).Update("dirt_level", 700).Error; err != nil {
		return err
	}

	return nil
}
