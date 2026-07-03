package services_test

import (
	"backend/models"
	"backend/services"
	"testing"
)

// TestCreateUserLifestyle は新規ユーザーへの生活環境情報登録を確認する
func TestCreateUserLifestyle(t *testing.T) {

	truncateUsersAndRooms(t) // usersテーブルとuser_roomsテーブルを初期化する

	req := services.LifestyleRequest{
		IsAlone:       true,  // 一人暮らし
		HasWasher:     false, // 洗濯機なし
		HasVacuum:     true,  // 掃除機あり
		HasRobot:      false, // ロボット掃除機なし
		UseTableware:  true,  // 食器を使用する
		HasDishwasher: false, // 食洗機なし
	}

	res, err := services.CreateUserLifestyle("user-001", req)
	if err != nil {
		t.Fatalf("CreateUserLifestyle failed: %v", err) // 登録処理自体が失敗したら終了
	}

	if !res.IsAlone || res.HasWasher || !res.HasVacuum || res.HasRobot || !res.UseTableware || res.HasDishwasher {
		t.Fatalf("unexpected response: %+v", res) // レスポンスの内容が期待通りでなければ失敗
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("UserRoom not created: %v", err) // DBに行が作成されていなければ終了
	}
	if !room.IsAlone || room.HasWasher || !room.HasVacuum || room.HasRobot || !room.UseTableware || room.HasDishwasher {
		t.Fatalf("unexpected UserRoom: %+v", room) // DBの値が期待通りでなければ失敗
	}
}

// TestCreateUserLifestyle_ExistingRoom は register 済みユーザーへの登録が UPSERT になることを確認する
func TestCreateUserLifestyle_ExistingRoom(t *testing.T) {

	truncateUsersAndRooms(t) // usersテーブルとuser_roomsテーブルを初期化する

	// register済みなど、既に UserRoom が存在するケースを再現
	existing := &models.UserRoom{
		UserID:        "user-001",
		IsAlone:       false, // 実家暮らし
		HasWasher:     true,  // 洗濯機あり
		HasVacuum:     true,  // 掃除機あり
		HasRobot:      true,  // ロボット掃除機あり
		UseTableware:  true,  // 食器を使用する
		HasDishwasher: true,  // 食洗機あり
	}
	if err := models.DB.Create(existing).Error; err != nil {
		t.Fatal(err) // 事前準備が失敗したら終了
	}

	req := services.LifestyleRequest{
		IsAlone:       true,  // 一人暮らしに変更
		HasWasher:     false, // 洗濯機なしに変更
		HasVacuum:     false, // 掃除機なしに変更
		HasRobot:      false, // ロボット掃除機なしに変更
		UseTableware:  false, // 食器を使用しないに変更
		HasDishwasher: false, // 食洗機なしに変更
	}

	if _, err := services.CreateUserLifestyle("user-001", req); err != nil {
		t.Fatalf("CreateUserLifestyle failed: %v", err) // UPSERT処理が失敗したら終了
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("UserRoom not found: %v", err) // 行が見つからなければ終了
	}
	if room.IsAlone != true || room.HasWasher || room.HasVacuum || room.HasRobot || room.UseTableware || room.HasDishwasher {
		t.Fatalf("unexpected UserRoom after upsert: %+v", room) // 全カラムが上書きされていなければ失敗
	}
}

// TestUpdateUserLifestyle は既存ユーザーの生活環境情報編集を確認する
func TestUpdateUserLifestyle(t *testing.T) {

	truncateUsersAndRooms(t) // usersテーブルとuser_roomsテーブルを初期化する

	existing := &models.UserRoom{
		UserID:        "user-001",
		IsAlone:       false, // 実家暮らし
		HasWasher:     true,  // 洗濯機あり
		HasVacuum:     true,  // 掃除機あり
		HasRobot:      true,  // ロボット掃除機あり
		UseTableware:  true,  // 食器を使用する
		HasDishwasher: true,  // 食洗機あり
	}
	if err := models.DB.Create(existing).Error; err != nil {
		t.Fatal(err) // 事前準備が失敗したら終了
	}

	req := services.LifestyleRequest{
		IsAlone:       true,  // 一人暮らしに変更
		HasWasher:     false, // 洗濯機なしに変更
		HasVacuum:     false, // 掃除機なしに変更
		HasRobot:      false, // ロボット掃除機なしに変更
		UseTableware:  false, // 食器を使用しないに変更
		HasDishwasher: false, // 食洗機なしに変更
	}

	if err := services.UpdateUserLifestyle("user-001", req); err != nil {
		t.Fatalf("UpdateUserLifestyle failed: %v", err) // 更新処理が失敗したら終了
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("UserRoom not found: %v", err) // 行が見つからなければ終了
	}
	if room.IsAlone != true || room.HasWasher || room.HasVacuum || room.HasRobot || room.UseTableware || room.HasDishwasher {
		t.Fatalf("unexpected UserRoom after update: %+v", room) // 全カラムが更新されていなければ失敗
	}
}
