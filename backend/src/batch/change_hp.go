package batch

import (
	"backend/models"
	"time"

	"gorm.io/gorm"
)

func ChangeHPTicker() {
		go func() {
		//3時間ごとにタスクを作成する
		ticker := time.NewTicker(3 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ChangeHP()
			}
		}
	}()
}


func ChangeHP() error {
	//汚さが72以下のユーザーのHPを60増加させる
	if err := models.DB.Model(&models.User{}).Where("dirt_level <= ?", 72).Update("health_point", gorm.Expr("health_point + ?", 60)).Error; err != nil {
		return err
	}
	//HPの上限は1000なので、HPが1000を超えたユーザーはHPを1000に固定する
	if err := models.DB.Model(&models.User{}).Where("health_point > ?", 1000).Update("health_point", 1000).Error; err != nil {
		return err
	}

	//汚さが73以上のユーザーのHPは(汚さ ÷  285) × 10　HPを減らす
	if err := models.DB.Model(&models.User{}).Where("dirt_level >= ?", 73).Update("health_point", gorm.Expr("health_point - ((dirt_level / 285) * 10)")).Error; err != nil {
		return err
	}

	return nil
}

