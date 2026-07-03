package services_test

import (
	"backend/models"
	"backend/services"
	"testing"
)

func TestCreateUserLifestyle(t *testing.T) {

	truncateUsersAndRooms(t)

	req := services.LifestyleRequest{
		IsAlone:   true,
		HasWasher: false,
		HasVacuum: true,
	}

	res, err := services.CreateUserLifestyle("user-001", req)
	if err != nil {
		t.Fatalf("CreateUserLifestyle failed: %v", err)
	}

	if !res.IsAlone || res.HasWasher || !res.HasVacuum {
		t.Fatalf("unexpected response: %+v", res)
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("UserRoom not created: %v", err)
	}
	if !room.IsAlone || room.HasWasher || !room.HasVacuum {
		t.Fatalf("unexpected UserRoom: %+v", room)
	}
}

func TestCreateUserLifestyle_ExistingRoom(t *testing.T) {

	truncateUsersAndRooms(t)

	// register済みなど、既に UserRoom が存在するケースを再現
	existing := &models.UserRoom{UserID: "user-001", IsAlone: false, HasWasher: true, HasVacuum: true}
	if err := models.DB.Create(existing).Error; err != nil {
		t.Fatal(err)
	}

	req := services.LifestyleRequest{
		IsAlone:   true,
		HasWasher: false,
		HasVacuum: false,
	}

	if _, err := services.CreateUserLifestyle("user-001", req); err != nil {
		t.Fatalf("CreateUserLifestyle failed: %v", err)
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("UserRoom not found: %v", err)
	}
	if !room.IsAlone || room.HasWasher || room.HasVacuum {
		t.Fatalf("unexpected UserRoom after upsert: %+v", room)
	}
}

func TestUpdateUserLifestyle(t *testing.T) {

	truncateUsersAndRooms(t)

	existing := &models.UserRoom{UserID: "user-001", IsAlone: false, HasWasher: true, HasVacuum: true}
	if err := models.DB.Create(existing).Error; err != nil {
		t.Fatal(err)
	}

	req := services.LifestyleRequest{
		IsAlone:   true,
		HasWasher: false,
		HasVacuum: false,
	}

	if err := services.UpdateUserLifestyle("user-001", req); err != nil {
		t.Fatalf("UpdateUserLifestyle failed: %v", err)
	}

	var room models.UserRoom
	if err := models.DB.First(&room, "user_id = ?", "user-001").Error; err != nil {
		t.Fatalf("UserRoom not found: %v", err)
	}
	if !room.IsAlone || room.HasWasher || room.HasVacuum {
		t.Fatalf("unexpected UserRoom after update: %+v", room)
	}
}
