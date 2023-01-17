package store

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	common2 "github.com/gensha256/data_collector/pkg/common"
	"github.com/gensha256/data_collector/pkg/models"

	"github.com/go-redis/redis/v8"
)

const (
	cmcSet  = "cmczet"
	expired = 24 * 60 * 60 * time.Second
)

type RedisStore struct {
	rds *redis.Client
}

func NewRedisStore() (*RedisStore, error) {

	conf := common2.NewConfig()

	rds := redis.NewClient(&redis.Options{
		Addr:     conf.RedisHost,
		Username: conf.RedisUsername,
		Password: conf.RedisPassword,
		DB:       0,
	})

	err := rds.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &RedisStore{rds: rds}, nil
}

func (rs *RedisStore) StoreCmcEntity(ctx context.Context, entity models.CmcEntity) error {
	key := entity.GetRedisKey()
	val := entity.GetAsJSON()

	err := rs.rds.Set(ctx, key, val, expired).Err()
	if err != nil {

		return err
	}

	err = rs.rds.ZAdd(ctx, cmcSet, &redis.Z{
		Score:  float64(entity.CmcRank),
		Member: entity.Symbol,
	}).Err()

	if err != nil {
		return err
	}

	return nil
}

func (rs *RedisStore) GetCmcEntityTimeSeriesBySymbol(ctx context.Context, symbol string) ([]models.CmcEntity, error) {
	keyPattern := fmt.Sprintf(models.CmcEntityBySymbolPattern, strings.ToLower(symbol))

	keysMatched := make([]string, 0)
	var cursor uint64
	var err error
	var keys []string

	for {

		keys, cursor, err = rs.rds.Scan(ctx, cursor, keyPattern, 100).Result()
		if err != nil {
			return nil, err
		}

		keysMatched = append(keysMatched, keys...)
		if cursor == 0 {
			break
		}
	}

	result := make([]models.CmcEntity, 0)

	for _, value := range keysMatched {
		cmcAsJSON, err := rs.rds.Get(ctx, value).Result()
		if err != nil {
			return nil, err
		}

		var entity models.CmcEntity
		err = json.Unmarshal([]byte(cmcAsJSON), &entity)
		if err != nil {
			return nil, err
		}

		result = append(result, entity)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].LastUpdated.Before(result[j].LastUpdated)
	})

	return result, nil
}

func (rs *RedisStore) GetSymbols(ctx context.Context) ([]string, error) {
	sortSymbols, err := rs.rds.ZRange(ctx, cmcSet, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return sortSymbols, err
}

func (rs *RedisStore) JustStore(entity models.CmcEntity) error {
	key := entity.GetRedisKey()
	val := entity.GetAsJSON()
	err := rs.rds.Set(context.Background(), key, val, expired).Err()
	if err != nil {
		return err
	}
	return nil
}
