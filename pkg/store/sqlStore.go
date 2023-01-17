package store

import (
	"context"
	"log"
	"strconv"

	"github.com/gensha256/data_collector/pkg/common"
	"github.com/gensha256/data_collector/pkg/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PgxStore struct {
	pgx *pgx.Conn
}

func NewPgxConnect() (*PgxStore, error) {
	conf := common.NewConfig()

	portInt, _ := strconv.Atoi(conf.PgxPort)
	port := uint16(portInt)

	config, _ := pgx.ParseConfig("")
	config.Host = conf.PgxHost
	config.Port = port
	config.User = conf.PgxUser
	config.Password = conf.PgxPassword
	config.Database = conf.PgxDbName

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Println("DB connection error", err)
	}

	return &PgxStore{pgx: conn}, nil
}

func (pg *PgxStore) CreateTableUsers() error {
	sqlCreateTable := `CREATE TABLE IF NOT EXISTS users
	(
		id VARCHAR(250) NOT NULL PRIMARY KEY,
		email VARCHAR(150),
    	telegram VARCHAR(50),
   		
    	UNIQUE (email, telegram)
    );`

	rows, err := pg.pgx.Exec(context.Background(), sqlCreateTable)
	if err != nil {
		return err
	}

	log.Println(rows)
	return nil
}

func (pg *PgxStore) CreateUser(us models.User) error {
	id := uuid.New()
	strId := id.String()

	res, err := pg.pgx.Exec(context.Background(),
		`INSERT INTO users (id,email,telegram) VALUES ($1, $2, $3)`,
		strId, us.Email, us.Telegram)

	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Println(res)

	return nil
}

func (pg *PgxStore) DeleteUser(id string) error {
	res, err := pg.pgx.Exec(context.Background(), `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	log.Println(res)

	return nil
}

func (pg *PgxStore) UpdateUser(us models.User) error {
	res, err := pg.pgx.Exec(context.Background(), `UPDATE users SET email=$1 WHERE id=$2`, us.Email, us.Id)
	if err != nil {
		return err
	}
	log.Println(res)

	return nil
}

func (pg *PgxStore) GetAllUsers() ([]models.User, error) {
	rows, err := pg.pgx.Query(context.Background(), `SELECT * FROM users`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	usr := models.User{}
	var res []models.User

	for rows.Next() {
		err := rows.Scan(&usr.Id, &usr.Email, &usr.Telegram)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		res = append(res, usr)
	}

	return res, nil
}

func (pg *PgxStore) GetUserId(us models.User) (string, error) {
	rows, err := pg.pgx.Query(context.Background(), `SELECT id FROM users WHERE email=$1`, us.Email)
	if err != nil {
		log.Fatal(err)
	}

	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
	}

	return id, nil
}
