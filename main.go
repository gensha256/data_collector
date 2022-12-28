package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gensha256/data_collector/cmc"
	"github.com/gensha256/data_collector/store"
)

func main() {
	rds, err := store.NewRedisStore()
	if err != nil {
		log.Fatal(err)
	}

	cmcAPI := cmc.NewAPI()

	//TODO: Remove this handler, then add cron job to fetch cmc crypto once per hour
	http.HandleFunc("/cmc", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		cmcData := cmcAPI.GetCryptoLatest()

		for _, value := range cmcData {

			err := rds.StoreCmcEntity(value)
			if err != nil {
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		byteArr, err := json.Marshal(cmcData)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(byteArr)
	})

	http.HandleFunc("/symbols", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		symbolsData, err := rds.GetSymbols()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		byteArr, err := json.Marshal(symbolsData)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(byteArr)
	})

	http.HandleFunc("/cmc/", func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		requestPath := req.URL.Path
		splitPath := strings.Split(requestPath, "/cmc/")

		//TODO: Validate path for create usage for insure split path 1
		dataBySymbol, err := rds.GetCmcEntityTimeSeriesBySymbol(splitPath[1])
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		byteArr, err := json.Marshal(dataBySymbol)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(byteArr)
	})

	log.Println("Listening on a 8080 port...")
	err = http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
