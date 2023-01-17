package tests

import (
	"context"
	"log"
	"testing"

	"github.com/gensha256/data_collector/pkg/common"
	"github.com/gensha256/data_collector/pkg/store"
)

func TestPgxConnect(t *testing.T) {
	conf := common.NewConfig()

	db := store.NewPgxConnect()
	if db.Config().User != conf.PgxUser {
		t.Error("Not correct user name ")
		t.Fail()
	}

	err := store.UpdateValue(db, "Max", "d2c1c48a-5e49-4dea-9ea3-8b8c162c7438")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	rows, err := db.Query(context.Background(), `SELECT name FROM users WHERE id = $id`)

	var name string

	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			t.Error(err)
		}
		if name != "Peter" {
			t.Error("Name not changed")
			t.Fail()
		}
	}
}

func TestDeleteValue(t *testing.T) {
	db := store.NewPgxConnect()

	err := store.DeleteValue(db, "735ce1ac-6a56-45e4-bd69-144afebb705f")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	rows, err := db.Query(context.Background(), `SELECT id FROM users`)
	if err != nil {
		t.Error(err)
	}

	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
		if id == "735ce1ac-6a56-45e4-bd69-144afebb705f" {
			t.Error("Not deleting id from users")
			t.Fail()
		}
	}
}

func TestInsertValue(t *testing.T) {
	db := store.NewPgxConnect()

	err := store.CreateTableUsers(db)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	err = store.InsertValue(db, "Som", "Dop", "Male", "somaal123@", "+38094545232")
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestGetAll(t *testing.T) {
	db := store.NewPgxConnect()

	res, err := store.GetAllValue(db)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	log.Println(res)
	if len(res) == 0 {
		t.Error("Not correct table length")
		t.Fail()
	}
}
