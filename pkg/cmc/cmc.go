package cmc

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gensha256/data_collector/pkg/common"
	common2 "github.com/gensha256/data_collector/pkg/models"
)

const (
	ApiKey = "X-CMC_PRO_API_KEY"
	ApiUrl = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?start=1&limit=%d"
)

type API struct {
	conf *common.Conf
}

func NewAPI() *API {
	conf := common.NewConfig()

	return &API{conf: conf}
}

func (c *API) GetCryptoLatest() []common2.CmcEntity {

	client := http.Client{}
	url := fmt.Sprintf(ApiUrl, c.conf.CmcApiLimit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add(ApiKey, c.conf.CmcApiToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Status :", res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var cmcResult common2.LatestResponse
	err = json.Unmarshal(body, &cmcResult)
	if err != nil {
		log.Fatal(err)
	}

	result := make([]common2.CmcEntity, 0)

	for _, val := range cmcResult.Data {
		entity := common2.CmcEntity{
			Id:               val.Id,
			Name:             val.Name,
			Symbol:           val.Symbol,
			MaxSupply:        val.MaxSupply,
			CmcRank:          val.CmcRank,
			Price:            val.Quote["USD"].Price,
			Volume24h:        val.Quote["USD"].Volume24h,
			VolumeChange24h:  val.Quote["USD"].VolumeChange24h,
			PercentChange1h:  val.Quote["USD"].PercentChange1h,
			PercentChange24h: val.Quote["USD"].PercentChange24h,
			MarketCap:        val.Quote["USD"].MarketCap,
			LastUpdated:      val.LastUpdated,
		}
		result = append(result, entity)
	}

	return result
}
