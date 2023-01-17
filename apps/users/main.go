package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gensha256/data_collector/pkg/models"
	"github.com/gensha256/data_collector/pkg/store"
)

func main() {

	db, err := store.NewPgxConnect()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/users", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if req.Method != http.MethodGet {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		resUsers, err := db.GetAllUsers()
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
	})

	http.HandleFunc("/insert", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if req.Method != http.MethodPost {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user := models.User{}

		byteArr, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatal()
		}

		err = json.Unmarshal(byteArr, &user)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		err = db.CreateUser(user)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
		writer.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/update", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if req.Method != http.MethodPut {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		us := models.User{}

		byteArr, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatal()
		}

		err = json.Unmarshal(byteArr, &us)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		err = db.UpdateUser(us)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
		writer.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/delete/", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if req.Method != http.MethodDelete {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		requestPath := req.URL.Path
		splitPath := strings.Split(requestPath, "/")
		id := splitPath[len(splitPath)-1]

		err := db.DeleteUser(id)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}
		writer.WriteHeader(http.StatusOK)
	})

	log.Println("Listening on a 8080 port...")
	http.ListenAndServe(":8080", nil)
}
