package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

const GlobalPathToConfig = "/home/gen/space/go/src/github.com/gensha256/data_collector/config.yaml"

type Conf struct {
	RedisHost     string `yaml:"redis_host"`
	RedisPassword string `yaml:"redis_password"`
	RedisUsername string `yaml:"redis_username"`
	CmcApiToken   string `yaml:"cmc_api_token"`
	CmcApiLimit   int    `yaml:"cmc_api_limit"`
	PgxHost       string `yaml:"pgx_host"`
	PgxPort       string `yaml:"pgx_port"`
	PgxUser       string `yaml:"pgx_user"`
	PgxPassword   string `yaml:"pgx_password"`
	PgxDbName     string `yaml:"pgx_db_name"`
}

func NewConfig() *Conf {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = GlobalPathToConfig
	}

	fileData, err := os.ReadFile(GlobalPathToConfig)
	if err != nil {
		panic(err)
	}

	c := &Conf{}
	err = yaml.Unmarshal(fileData, c)
	if err != nil {
		panic(err)
	}

	return c
}
