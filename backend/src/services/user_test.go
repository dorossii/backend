package services_test

import (
	"backend/models"
	"backend/services"
	"errors"
	"log"
	"testing"

	"gorm.io/gorm"
)

func truncateUsersAndRooms(t *testing.T) {
	t.Helper()

	if err := models.DB.Exec("TRUNCATE TABLE user_rooms").Error; err != nil {
		t.Fatal(err)
	}
	if err := models.DB.Exec("TRUNCATE TABLE users").Error; err != nil {
		t.Fatal(err)
	}
}

func TestRegisterUser(t *testing.T) {

	truncateUsersAndRooms(t)

	req := services.RegisterUserRequest{
		BirthDate:  946684800, // 2000-01-01 00:00:00 UTC
		LivingType: "alone",
	}

	res, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.UserID != "user-001" {
		t.Fatalf("unexpected UserID: %s", res.UserID)
	}
	if res.UserName != "syatyo" {
		t.Fatalf("unexpected UserName: %s", res.UserName)
	}
	if res.BirthDate != 946684800 {
		t.Fatalf("unexpected BirthDate: %d", res.BirthDate)
	}
	if res.LivingType != "alone" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("UserRoom not created: %v", err)
	}
	if !room.IsAlone {
		t.Fatalf("expected IsAlone=true for livingType=alone")
	}
	if !room.HasWasher || !room.HasVacuum || !room.HasRobot || !room.UseTableware || !room.HasDishwasher {
		t.Fatalf("expected default lifestyle fields to be true: %+v", room) // register時点では全てデフォルト値(true)であることを確認する
	}

	var user models.User
	if err := models.DB.First(&user, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("User not found: %v", err)
	}

	log.Printf("User icon: %s, bgColor: %s", user.Icon, user.BgColor)

	// ユーザー二人目を追加
	req = services.RegisterUserRequest{
		BirthDate:  926726400, // 1999-05-15 00:00:00 UTC
		LivingType: "alone",
	}

	res, err = services.RegisterUser("user-002", "goro", "goro@example.com", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.UserID != "user-002" {
		t.Fatalf("unexpected UserID: %s", res.UserID)
	}
	if res.UserName != "goro" {
		t.Fatalf("unexpected UserName: %s", res.UserName)
	}
	if res.BirthDate != 926726400 {
		t.Fatalf("unexpected BirthDate: %d", res.BirthDate)
	}
	if res.LivingType != "alone" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}

	// ユーザー003
	req = services.RegisterUserRequest{
		BirthDate:  946684800, // 2000-01-01 00:00:00 UTC
		LivingType: "alone",
	}

	res, err = services.RegisterUser("user-003", "usr003", "usr003@example.com", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.UserID != "user-003" {
		t.Fatalf("unexpected UserID: %s", res.UserID)
	}

	if res.UserName != "usr003" {
		t.Fatalf("unexpected UserName: %s", res.UserName)
	}

	if res.BirthDate != 946684800 {
		t.Fatalf("unexpected BirthDate: %d", res.BirthDate)
	}

	if res.LivingType != "alone" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}
}

func TestRegisterUser_Family(t *testing.T) {

	truncateUsersAndRooms(t)

	req := services.RegisterUserRequest{
		BirthDate:  926726400, // 1999-05-15 00:00:00 UTC
		LivingType: "family",
	}

	res, err := services.RegisterUser("user-002", "goro", "goro@example.com", req)
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	if res.LivingType != "family" {
		t.Fatalf("unexpected LivingType: %s", res.LivingType)
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-002").Error; err != nil {
		t.Fatalf("UserRoom not created: %v", err)
	}
	if room.IsAlone {
		t.Fatalf("expected IsAlone=false for livingType=family")
	}
}

func TestRegisterUser_InvalidLivingType(t *testing.T) {

	truncateUsersAndRooms(t)

	req := services.RegisterUserRequest{
		BirthDate:  946684800,
		LivingType: "invalid",
	}

	_, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", req)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestGetUserStatus(t *testing.T) {

	truncateUsersAndRooms(t)

	req := services.RegisterUserRequest{
		BirthDate:  946684800,
		LivingType: "alone",
	}

	if _, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", req); err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	res, err := services.GetUserStatus("user-001")
	if err != nil {
		t.Fatalf("GetUserStatus failed: %v", err)
	}

	if res.UserID != "user-001" {
		t.Fatalf("unexpected UserID: %s", res.UserID)
	}
	if res.UserName != "syatyo" {
		t.Fatalf("unexpected UserName: %s", res.UserName)
	}
	if res.DirtLevel != 1 {
		t.Fatalf("unexpected DirtLevel: %d", res.DirtLevel)
	}
	if res.HealthPoint != 1000 {
		t.Fatalf("unexpected HealthPoint: %d", res.HealthPoint)
	}
}


func TestGetUserStatus_NotFound(t *testing.T) {

	truncateUsersAndRooms(t)

	_, err := services.GetUserStatus("not-exist-user")
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got: %v", err)
	}
}

// TestUpdateUserSetting はユーザー名が正しく更新されることを確認する
func TestUpdateUserSetting(t *testing.T) {

	truncateUsersAndRooms(t) // usersテーブルとuser_roomsテーブルを初期化する

	registerReq := services.RegisterUserRequest{
		BirthDate:  946684800, // 事前にユーザーを登録するための誕生日
		LivingType: "alone",   // 事前にユーザーを登録するための生活形態
	}
	if _, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", registerReq); err != nil {
		t.Fatalf("RegisterUser failed: %v", err) // 事前準備のユーザー登録が失敗したら終了
	}

	icon := "pineTree"
	bgColor := "icon1"
	
	if err := services.UpdateUserSetting("user-001", services.UserSettingRequest{
		UserName: "new-name",
		Icon:     &icon,
		BgColor:  &bgColor,
	}); err != nil {
		t.Fatalf("UpdateUserSetting failed: %v", err) // 更新処理自体が失敗したら終了
	}

	var user models.User
	if err := models.DB.First(&user, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("failed to find User: %v", err) // 更新後のユーザーが取得できなければ終了
	}
	if user.UserName != "new-name" {
		t.Fatalf("unexpected UserName: %s", user.UserName) // ユーザー名が更新されていなければ失敗
	}
	// icon/backgroundを指定しない場合はRegisterUser時のデフォルト値が維持される
	if user.Icon != "pineTree" {
		t.Fatalf("unexpected Icon: %s", user.Icon)
	}
	if user.BgColor != "icon1" {
		t.Fatalf("unexpected BgColor: %s", user.BgColor)
	}
}

// TestUpdateUserSetting_IconAndBackground はicon/backgroundが指定通り更新されることを確認する
func TestUpdateUserSetting_IconAndBackground(t *testing.T) {

	truncateUsersAndRooms(t)

	registerReq := services.RegisterUserRequest{
		BirthDate:  946684800,
		LivingType: "alone",
	}
	if _, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", registerReq); err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	icon := "dog"
	bgColor := "icon2"

	req := services.UserSettingRequest{
		UserName: "new-name",
		Icon:     &icon,
		BgColor:  &bgColor,
	}
	if err := services.UpdateUserSetting("user-001", req); err != nil {
		t.Fatalf("UpdateUserSetting failed: %v", err)
	}
 
	var user models.User
	if err := models.DB.First(&user, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("failed to find User: %v", err)
	}
	if user.Icon != "dog" {
		t.Fatalf("unexpected Icon: %s", user.Icon)
	}
	if user.BgColor != "icon2" {
		t.Fatalf("unexpected BgColor: %s", user.BgColor)
	}
}

// TestUpdateUserSetting_EmptyUserName は空文字のユーザー名がエラーになることを確認する
func TestUpdateUserSetting_EmptyUserName(t *testing.T) {

	truncateUsersAndRooms(t) // usersテーブルとuser_roomsテーブルを初期化する

	registerReq := services.RegisterUserRequest{
		BirthDate:  946684800, // 事前にユーザーを登録するための誕生日
		LivingType: "alone",   // 事前にユーザーを登録するための生活形態
	}
	if _, err := services.RegisterUser("user-001", "syatyo", "syatyo@example.com", registerReq); err != nil {
		t.Fatalf("RegisterUser failed: %v", err) // 事前準備のユーザー登録が失敗したら終了
	}

	err := services.UpdateUserSetting("user-001", services.UserSettingRequest{UserName: ""}) // 空文字を渡す
	if err == nil {
		t.Fatal("expected error but got nil") // バリデーションエラーが返らなければ失敗
	}
}
