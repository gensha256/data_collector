package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gensha256/data_collector/pkg/store"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Telegram string `json:"telegram"`
}

func main() {

	http.HandleFunc("/users", getUsers)
	http.HandleFunc("/insert", postUsers)
	http.HandleFunc("/update", putUsers)
	http.HandleFunc("/delete/", deleteUsers)

	log.Println("Listening on a 8080 port...")
	http.ListenAndServe(":8080", nil)
}

func getUsers(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	db := store.NewPgxConnect()

	if req.Method != http.MethodGet {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resUsers, err := store.GetAllValue(db)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	byteArr, err := json.Marshal(resUsers)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = writer.Write(byteArr)
}

func postUsers(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	db := store.NewPgxConnect()

	if req.Method != http.MethodPost {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := &User{}
	err := json.NewDecoder(req.Body).Decode(user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	err = store.InsertValue(db, user.Name, user.Surname, user.Gender, user.Email, user.Telegram)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	writer.WriteHeader(http.StatusOK)

	defer db.Close(context.Background())
}

func putUsers(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	db := store.NewPgxConnect()

	if req.Method != http.MethodPut {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	us := &User{}
	err := json.NewDecoder(req.Body).Decode(us)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	err = store.UpdateValue(db, us.Name, us.Id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	writer.WriteHeader(http.StatusOK)

	defer db.Close(context.Background())
}

func deleteUsers(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	db := store.NewPgxConnect()

	if req.Method != http.MethodDelete {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	requestPath := req.URL.Path
	splitPath := strings.Split(requestPath, "/")
	id := splitPath[len(splitPath)-1]

	err := store.DeleteValue(db, id)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}
	writer.WriteHeader(http.StatusOK)

	defer db.Close(context.Background())
}
