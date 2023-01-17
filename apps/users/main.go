package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gensha256/data_collector/pkg/models"
	"github.com/gensha256/data_collector/pkg/store"
)

func main() {

	userDB := store.NewUserStore()

	http.HandleFunc("/users", func(writer http.ResponseWriter, req *http.Request) {

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		method := req.Method

		switch method {
		case http.MethodGet:

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

		case http.MethodPost:

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

			_, _ = writer.Write(byteArr)

		case http.MethodPut:

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

			_, _ = writer.Write(byteArr)

		default:
			http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
