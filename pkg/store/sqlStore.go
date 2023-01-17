package store

import (
	"context"
	"github.com/google/uuid"
	"log"
	"strconv"

	"github.com/gensha256/data_collector/pkg/common"

	"github.com/jackc/pgx/v5"
)

type users struct {
	Id       string
	Name     string
	SurName  string
	Gender   string
	Email    string
	Telegram string
}

func NewPgxConnect() *pgx.Conn {
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

	return conn
}

func CreateTableUsers(pg *pgx.Conn) error {
	sqlCreateTable := `CREATE TABLE IF NOT EXISTS users
	(
		id VARCHAR(250) NOT NULL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		surname VARCHAR(50) NOT NULL,
		gender VARCHAR(7) NOT NULL,
		email VARCHAR(150),
    	telegram VARCHAR(50) 
    );`

	rows, err := pg.Exec(context.Background(), sqlCreateTable)
	if err != nil {
		return err
	}

	log.Println(rows)
	return nil
}

func InsertValue(pg *pgx.Conn, nm string, sur string, gen string, email string, tg string) error {
	id := uuid.New()
	strId := id.String()

	res, err := pg.Exec(context.Background(), `INSERT INTO users (id, name, surname, gender, email, telegram) VALUES ($1, $2, $3, $4, $5, $6)`, strId, nm, sur, gen, email, tg)
	if err != nil {
		return err
	}
	log.Println(res)

	return nil
}

func DeleteValue(pg *pgx.Conn, id string) error {
	res, err := pg.Exec(context.Background(), `DELETE FROM users WHERE  id = $1`, id)
	if err != nil {
		return err
	}
	log.Println(res)

	return nil
}

func UpdateValue(pg *pgx.Conn, name string, id string) error {
	res, err := pg.Exec(context.Background(), `UPDATE users SET name=$1 WHERE id=$2`, name, id)
	if err != nil {
		return err
	}
	log.Println(res)

	return nil
}

func GetAllValue(pg *pgx.Conn) ([]users, error) {
	rows, err := pg.Query(context.Background(), `SELECT * FROM users`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	usr := users{}
	var res []users

	for rows.Next() {
		err := rows.Scan(&usr.Id, &usr.Name, &usr.SurName, &usr.Gender, &usr.Email, &usr.Telegram)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		res = append(res, usr)
	}

	return res, nil
}
