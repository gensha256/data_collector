package main

import (
	"context"
	"encoding/json"
	"github.com/gensha256/data_collector/pkg/cmc"
	"github.com/gensha256/data_collector/pkg/store"
	"log"
	"net/http"
	"strings"

	"github.com/robfig/cron/v3"
)

func main() {
	rds, err := store.NewRedisStore()
	if err != nil {
		log.Fatal(err)
	}

	cmcAPI := cmc.NewAPI()
	initCron(cmcAPI, rds)

	http.HandleFunc("/symbols", func(writer http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		symbolsData, err := rds.GetSymbols(ctx)
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

		ctx := req.Context()

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)

		requestPath := req.URL.Path
		splitPath := strings.Split(requestPath, "/")
		lastPath := splitPath[len(splitPath)-1]

		dataBySymbol, err := rds.GetCmcEntityTimeSeriesBySymbol(ctx, lastPath)
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

func initCron(api *cmc.API, rds *store.RedisStore) {

	crn := cron.New()

	_, err := crn.AddFunc("@hourly", func() {

		cmcData := api.GetCryptoLatest()

		for _, value := range cmcData {

			err := rds.StoreCmcEntity(context.Background(), value)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println("cmc cache :", len(cmcData))
	})

	if err != nil {
		log.Fatal(err)
	}
	crn.Start()
}
