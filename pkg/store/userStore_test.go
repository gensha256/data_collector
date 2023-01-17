package store

import (
	"testing"

	"github.com/gensha256/data_collector/pkg/models"
)

func TestUserStore(t *testing.T) {

	// Test data

	testEmail := "test@gmail.com"
	testEmail2 := "test2@gmail.com"
	testTelegram := "@test"
	testTelegram2 := "@test2"

	db := NewUserStore()

	// Make sure test data cleaned up
	_ = db.DeleteByEmail(testEmail)
	_ = db.DeleteByEmail(testEmail2)

	usr := models.User{
		Email:    testEmail,
		Telegram: testTelegram,
	}

	resultCreate, err := db.Create(usr)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if resultCreate.ID == "" {
		t.Error("not correct id")
		t.Fail()
	}

	if resultCreate.Email != usr.Email {
		t.Error("not correct email")
		t.Fail()
	}

	if resultCreate.Telegram != usr.Telegram {
		t.Error("not correct telegram")
		t.Fail()
	}

	userByID, err := db.GetById(resultCreate.ID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if userByID.ID != resultCreate.ID {
		t.Error("not correct id")
		t.Fail()
	}

	if userByID.Email != resultCreate.Email {
		t.Error("not correct email")
		t.Fail()
	}

	if userByID.Telegram != resultCreate.Telegram {
		t.Error("not correct telegram")
		t.Fail()
	}

	userByID.Email = testEmail2
	userByID.Telegram = testTelegram2

	updateUser, err := db.Update(userByID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if userByID.ID != updateUser.ID {
		t.Error("not correct id")
		t.Fail()
	}

	if userByID.Email != updateUser.Email {
		t.Error("not correct email")
		t.Fail()
	}

	if userByID.Telegram != updateUser.Telegram {
		t.Error("not correct telegram")
		t.Fail()
	}

	err = db.Delete(userByID.ID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
