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

const (
	sqlCreateTable = `CREATE TABLE IF NOT EXISTS users
	(
		id VARCHAR(250) NOT NULL PRIMARY KEY,
		email VARCHAR(150),
    	telegram VARCHAR(50),
   		
    	UNIQUE (email, telegram)
    );`

	sqlCreate = `INSERT INTO users (id, email, telegram) VALUES ($1, $2, $3);`
	sqlUpdate = `UPDATE users SET email=$1, telegram=$2 WHERE id=$3;`
	sqlDelete = `DELETE FROM users WHERE id=$1;`

	sqlSelectAll        = `SELECT id, email, telegram FROM users;`
	sqlSelectByID       = `SELECT id, email, telegram FROM users WHERE id=$1;`
	sqlSelectByEmail    = `SELECT id, email, telegram FROM users WHERE email=$1;`
	sqlSelectByTelegram = `SELECT id, email, telegram FROM users WHERE telegram=$1;`
)

type UserStore struct {
	pgx *pgx.Conn
}

func NewUserStore() *UserStore {
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
		log.Fatal("DB connection error", err)
	}

	store := UserStore{pgx: conn}
	err = store.CreateTable()
	if err != nil {
		log.Fatal("error create user table", err)
	}

	return &store
}

func (pg *UserStore) CreateTable() error {
	rows, err := pg.pgx.Exec(context.Background(), sqlCreateTable)
	if err != nil {
		return err
	}

	log.Println(rows)
	return nil
}

func (pg *UserStore) Create(user models.User) (models.User, error) {
	id := uuid.New().String()

	result := models.User{}

	_, err := pg.pgx.Exec(
		context.Background(),
		sqlCreate,
		id,
		user.Email,
		user.Telegram)

	if err != nil {
		log.Println(err)
		return result, err
	}

	return pg.GetById(id)
}

func (pg *UserStore) Update(user models.User) (models.User, error) {
	_, err := pg.pgx.Exec(
		context.Background(),
		sqlUpdate,
		user.Email,
		user.Telegram,
		user.Id)
	if err != nil {
		return user, err
	}

	return pg.GetById(user.Id)
}

func (pg *UserStore) Delete(id string) error {
	_, err := pg.pgx.Exec(
		context.Background(),
		sqlDelete,
		id)
	if err != nil {
		return err
	}

	return nil
}

func (pg *UserStore) GetAll() ([]models.User, error) {
	rows, err := pg.pgx.Query(
		context.Background(),
		sqlSelectAll)

	if err != nil {
		log.Println("error on selectAll", err)
		return nil, err
	}

	result := make([]models.User, 0)

	for rows.Next() {
		usr := models.User{}
		err = rows.Scan(&usr.Id, &usr.Email, &usr.Telegram)
		if err != nil {
			log.Println("error on scan", err)
			return nil, err
		}

		result = append(result, usr)
	}

	return result, nil
}

func (pg *UserStore) GetById(id string) (models.User, error) {
	return pg.getBy(id, sqlSelectByID)
}

func (pg *UserStore) GetByEmail(email string) (models.User, error) {
	return pg.getBy(email, sqlSelectByEmail)
}

func (pg *UserStore) GetByTelegram(telegram string) (models.User, error) {
	return pg.getBy(telegram, sqlSelectByTelegram)
}

func (pg *UserStore) getBy(by string, sql string) (models.User, error) {
	row := pg.pgx.QueryRow(
		context.Background(),
		sql,
		by)

	user := models.User{}

	err := row.Scan(&user.Id, &user.Email, &user.Telegram)
	if err != nil {
		log.Println("error on user scan", err)
		return user, err
	}

	return user, nil
}
