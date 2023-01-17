package store

import (
	"testing"

	"github.com/gensha256/data_collector/pkg/models"
)

func TestPgxConnect(t *testing.T) {
	db, err := NewPgxConnect()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	err = db.CreateTableUsers()
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	usr := models.User{
		Email:    "test@gmail.com",
		Telegram: "@djindos0a4",
	}

	err = db.CreateUser(usr)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	userId, err := db.GetUserId(usr)
	if err != nil {
		t.Fail()
	}

	upUser := models.User{
		Id:    userId,
		Email: "test2@gmail.com",
	}

	err = db.UpdateUser(upUser)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	err = db.DeleteUser(userId)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
