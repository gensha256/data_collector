package store

import (
	"context"
	"testing"
	"time"

	"github.com/gensha256/data_collector/common"
)

var entity = common.CmcEntity{
	Id:               -1,
	Name:             "Test",
	Symbol:           "TEST",
	MaxSupply:        1,
	CmcRank:          -1,
	Price:            1,
	Volume24h:        1,
	VolumeChange24h:  -1,
	PercentChange1h:  1,
	PercentChange24h: 1,
	MarketCap:        -1,
	LastUpdated:      time.Date(2022, 12, 26, 13, 54, 00, 00, time.UTC),
}
var entityKey = "cmc.test.1672059600"
var entityTTS = int64(1672059600)

func TestRedisStore(t *testing.T) {
	_, err := NewRedisStore()
	if err != nil {
		t.Fail()
	}
}

func TestGetRedisKey(t *testing.T) {
	key := entity.GetRedisKey()
	if key != entityKey {
		t.Error("not correct redis key")
		t.Fail()
	}
}

func TestCmcEntityEvalTTS(t *testing.T) {
	result := entity.EvalTTS()

	if result != entityTTS {
		t.Error("not correct time")
		t.Fail()
	}
}

func TestStoreAndGetCmcEntity(t *testing.T) {
	rs, _ := NewRedisStore()

	ctx := context.Background()

	err := rs.StoreCmcEntity(ctx, entity)

	bySymbol, err := rs.GetCmcEntityTimeSeriesBySymbol(ctx, entity.Symbol)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	if len(bySymbol) == 0 {
		t.Error("invalid size of entities after store")
		t.Fail()
	}

	symbols, err := rs.GetSymbols(ctx)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	result := false

	for i := 0; i < len(symbols); i++ {
		if symbols[i] == entity.Symbol {
			result = true
			break
		}
	}

	if !result {
		t.Fail()
	}
}
