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

	userDB := store.NewUserStore()

	http.HandleFunc("/users", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if req.Method != http.MethodGet {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		usersList, err := userDB.GetAll()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		byteArr, err := json.Marshal(usersList)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(byteArr)
	})

	http.HandleFunc("/users/create", func(writer http.ResponseWriter, req *http.Request) {
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

		result, err := userDB.Create(user)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		byteArr, _ = json.Marshal(result)

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(byteArr)
	})

	http.HandleFunc("/users/update", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if req.Method != http.MethodPut {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user := models.User{}

		byteArr, err := io.ReadAll(req.Body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(byteArr, &user)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		user, err = userDB.Update(user)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		byteArr, _ = json.Marshal(user)

		writer.WriteHeader(http.StatusOK)
		_, _ = writer.Write(byteArr)
	})

	http.HandleFunc("/users/delete/", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		if req.Method != http.MethodDelete {
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		requestPath := req.URL.Path
		splitPath := strings.Split(requestPath, "/")
		id := splitPath[len(splitPath)-2]

		err := userDB.Delete(id)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusMethodNotAllowed)
			return
		}

		writer.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
