package cmc

import "time"

type LatestResponse struct {
	Data []ObjectResponse `json:"data"`
}

type ObjectResponse struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	MaxSupply   float64   `json:"max_supply"`
	CmcRank     int64     `json:"cmc_rank"`
	LastUpdated time.Time `json:"last_updated"`
	Quote       map[string]struct {
		Price            float64 `json:"price"`
		Volume24h        float64 `json:"volume_24h"`
		VolumeChange24h  float64 `json:"volume_change_24h"`
		PercentChange1h  float64 `json:"percent_change_1h"`
		PercentChange24h float64 `json:"percent_change_24h"`
		MarketCap        float64 `json:"market_cap"`
	} `json:"quote"`
}
