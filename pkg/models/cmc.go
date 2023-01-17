package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type CmcEntity struct {
	Id               int       `json:"id"`
	Name             string    `json:"name"`
	Symbol           string    `json:"symbol"`
	MaxSupply        float64   `json:"max_supply"`
	CmcRank          int64     `json:"cmc_rank"`
	Price            float64   `json:"price"`
	Volume24h        float64   `json:"volume_24_h"`
	VolumeChange24h  float64   `json:"volume_change_24h"`
	PercentChange1h  float64   `json:"percent_change_1h"`
	PercentChange24h float64   `json:"percent_change_24h"`
	MarketCap        float64   `json:"market_cap"`
	LastUpdated      time.Time `json:"last_updated"`
}

const (
	CmcEntityRedisPrefix     = "cmc"
	CmcEntityBySymbolPattern = "cmc.%s.*"
)

func (entity *CmcEntity) EvalTTS() int64 {
	date := entity.LastUpdated
	date = time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), 00, 00, 00, time.UTC)
	return date.Unix()
}

func (entity *CmcEntity) GetRedisKey() string {
	// cmc.btc.1671774245
	return strings.ToLower(
		fmt.Sprintf(
			"%s.%s.%d",
			CmcEntityRedisPrefix,
			entity.Symbol,
			entity.EvalTTS()))
}

func (entity *CmcEntity) GetAsJSON() string {
	byteArr, err := json.Marshal(entity)
	if err != nil {
		log.Fatal(err)
	}

	return string(byteArr)
}
