package store

import (
	"testing"

	"github.com/gensha256/data_collector/pkg/models"
)

func TestUserStore(t *testing.T) {

	db := NewUserStore()

	usr := models.User{
		Email:    "test@gmail.com",
		Telegram: "@test",
	}

	resultCreate, err := db.Create(usr)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if resultCreate.Id == "" {
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

	userByID, err := db.GetById(resultCreate.Id)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if userByID.Id != resultCreate.Id {
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

	userByID.Email = "test2@gmail.com"
	userByID.Telegram = "@test2"

	updateUser, err := db.Update(userByID)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if userByID.Id != updateUser.Id {
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

	err = db.Delete(userByID.Id)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
