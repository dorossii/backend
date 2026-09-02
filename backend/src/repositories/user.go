package repositories

import (
	"backend/models"
	"backend/utils"
	"errors"
	"time"
	"math/rand"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return models.DB.Create(user).Error
}

func CreateTaskForUser(userID string) error {

	// 1ユーザーにつき何個タスクを作成するかの定数
	const tasksPerUser = 2

	// 全てのBaseTaskを取得
	var baseTasks []models.BaseTask
	if err := models.DB.Find(&baseTasks).Error; err != nil {
		return err
	}

	// ベースタスクが2つ未満だと「1ユーザーにつき2タスク」を満たせないためエラーハンドリング
	if len(baseTasks) < tasksPerUser {
		return errors.New("insufficient base tasks available (minimum 2 required)")
	}

	// 乱数生成器の初期化
	r := rand.New(rand.NewSource(utils.NowJST().UnixNano()))

	var tasksToInsert []models.Task
	now := utils.NowJST()

	// 同じユーザーに同じタスクが重複して割り当たらないようにインデックスをシャッフル
	shuffledIndices := r.Perm(len(baseTasks))

	// 上位2つのランダムなタスクを選択
	for i := 0; i < tasksPerUser; i++ {
		baseTask := baseTasks[shuffledIndices[i]]

		// DueTime（期限）の仕様に合わせて終了時間を計算
		endTime := now.Add(time.Duration(baseTask.DueTime) * 24 * time.Hour)

		// UUIDを生成
		uuid, err := utils.Genid()
		if err != nil {
			return err
		}

		// imageflagがtrueの要素の中で10%の確率でRequireImageをtrueにする
		requireImage := false
		if baseTask.ImageFlag && r.Float64() < 0.1 {
			requireImage = true
		}

		task := models.Task{
			TaskID:       uuid,
			BaseID:       baseTask.BaseID,
			UserID:       userID,
			Status:       models.TaskStatusIncomplete,
			StartTime:    now,
			EndTime:      endTime,
			ImageID:      "", // 初期状態は空
			RequireImage: requireImage,
		}
		tasksToInsert = append(tasksToInsert, task)
	}

	// バルクインサート（一括保存）を実行
	if err := models.DB.CreateInBatches(&tasksToInsert, 100).Error; err != nil {
		return err
	}

	return nil
}


func CreateUserTx(tx *gorm.DB, user *models.User) error {
	return tx.Create(user).Error
}

// userRoomColumns は UserRoom の値を UPSERT/更新用の map に変換する
// bool フィールドに gorm:"default:true" が付与されているため、構造体経由の Create/Save だと
// ゼロ値(false)が「未設定」とみなされデフォルト値で上書きされてしまう。
// map 経由で書き込むことでゼロ値も確実に反映する。
func userRoomColumns(userID string, room models.UserRoom) map[string]interface{} {
	return map[string]interface{}{
		"user_id":        userID,             // ユーザーID
		"is_alone":       room.IsAlone,       // 一人暮らしか
		"has_washer":     room.HasWasher,     // 洗濯機があるか
		"has_vacuum":     room.HasVacuum,     // 掃除機があるか
		"has_robot":      room.HasRobot,      // ロボット掃除機があるか
		"use_tableware":  room.UseTableware,  // 食器を使用するか
		"has_dishwasher": room.HasDishwasher, // 食洗機があるか
	}
}

// CreateUserRoomTx は UserRoom を作成する
func CreateUserRoomTx(tx *gorm.DB, userRoom *models.UserRoom) error {
	return tx.
		Model(&models.UserRoom{}).                                // 対象テーブルを指定する
		Create(userRoomColumns(userRoom.UserID, *userRoom)).Error // map経由で全カラムを明示的に書き込む
}

// UpsertUserRoom は UserRoom を UPSERT する（存在すれば全カラム更新、なければ作成）
func UpsertUserRoom(userID string, room models.UserRoom) error {
	var count int64
	if err := models.DB.Model(&models.UserRoom{}).Where("user_id = ?", userID).Count(&count).Error; err != nil { // 既存行の有無を確認する
		return err
	}

	if count == 0 { // 存在しない場合は新規作成する
		return models.DB.
			Model(&models.UserRoom{}).
			Create(userRoomColumns(userID, room)).Error
	}

	return UpdateUserRoom(userID, room) // 存在する場合は更新する
}

// UpdateUserRoom は UserRoom の全カラムを更新する
func UpdateUserRoom(userID string, room models.UserRoom) error {
	return models.DB.
		Model(&models.UserRoom{}).
		Where("user_id = ?", userID). // 更新対象を user_id で絞り込む
		Updates(userRoomColumns(userID, room)).Error
}

func GetUser(UserID string) (*models.User, error) {
	var user models.User
	err := models.DB.First(&user, "user_id = ?", UserID).Error
	return &user, err
}

// GetUserRoom は UserRoom を取得する
func GetUserRoom(userID string) (*models.UserRoom, error) {
	var room models.UserRoom
	err := models.DB.First(&room, "user_id = ?", userID).Error
	return &room, err
}

// UpdateUserSetting は指定したユーザーの表示名・アイコン・背景を更新する
// icon/bgColor が nil の場合はそのカラムを更新しない
func UpdateUserSetting(userID string, userName string, icon *string, bgColor *string) error {
	columns := map[string]interface{}{
		"user_name": userName,
	}
	if icon != nil {
		columns["icon"] = *icon
	}
	if bgColor != nil {
		columns["bg_color"] = *bgColor
	}

	return models.DB.
		Model(&models.User{}).        // Userテーブルを対象にする
		Where("user_id = ?", userID). // 更新対象を user_id で絞り込む
		Updates(columns).Error
}

func UpdateAttackerSettings(userID string, targetUser string) error {
	err := models.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("target_user", targetUser).Error
	if err != nil {
		return err
	}

	return nil
}

func UpdateAttackerSettingsTx(tx *gorm.DB, userID string, targetUser string) error {
	err := tx.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("target_user", targetUser).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAttackerSettings は userID の target_user を返す	
func GetAttackerSettings(userID string) (string, error) {
	var user models.User
	err := models.DB.Select("target_user").First(&user, "user_id = ?", userID).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil // 設定がない場合は空文字を返す
	}
	if err != nil {
		return "", err
	}
	
	return user.TargetUser, nil
}

func UpdateDirtLevel(tx *gorm.DB, userID string, diff int) error {
	var user models.User

	err := tx.Model(&models.User{}).First(&user, "user_id = ?", userID).Error
	if err != nil {
		return err
	}

	// 計算後の値が0以下700以上ならぱわーでねじ曲げる
	newDirt := user.DirtLevel + diff

	if newDirt < 0 {
		newDirt = 0
	}

	if newDirt > 700 {
		newDirt = 700
	}

	return tx.Model(&models.User{}).Where("user_id = ?", userID).Update("dirt_level", newDirt).Error
}

// ListUsers は search が空でなければ UserName/Mailadress の部分一致で絞り込んだユーザー一覧を返す
func ListUsers(search string) ([]models.User, error) {
	var users []models.User
	query := models.DB.Model(&models.User{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("user_name LIKE ? OR mailadress LIKE ?", like, like)
	}
	err := query.Find(&users).Error
	return users, err
}

// UpdateUserName はユーザー名のみを更新する
func UpdateUserName(userID string, userName string) error {
	return models.DB.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("user_name", userName).Error
}

// UpdateUserStats はHealthPoint/DirtLevelを直接更新する
func UpdateUserStats(userID string, healthPoint *int, dirtLevel *int) error {
	updates := map[string]interface{}{}
	if healthPoint != nil {
		updates["health_point"] = *healthPoint
	}
	if dirtLevel != nil {
		updates["dirt_level"] = *dirtLevel
	}
	if len(updates) == 0 {
		return nil
	}

	return models.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updates).Error
}

