package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/godror/godror"
	"log"
	"os"
	"time"
)

type Repository struct {
	db *sql.DB
}

const (
	usersTableName    string = "users"
	groupsTableName   string = "groups"
	messagesTableName string = "messages"
)

func connectDb() *sql.DB {
	connString := os.Getenv("DB_CON_STRING")
	log.Printf(connString)
	db, err := sql.Open("godror", connString)
	if err != nil {
		panic(fmt.Errorf("error openning db connection %w", err))
	}
	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error pinging db: %w", err))
	}
	return db
}

func (r Repository) GetUsers() ([]User, error) {
	sqlQuery := fmt.Sprintf("SELECT id, tg_id, first_name, last_name, username, lang_code, is_bot, created FROM %s", usersTableName)
	rows, err := r.db.Query(sqlQuery)
	var users []User
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.Id,
			&user.TgId,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.LanguageCode,
			&user.IsBot,
			&user.Username,
			&user.Created,
		); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error on getting users %w", err)
	}
	return users, nil
}

func (r Repository) CreateUser(user User) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO %s (tg_id, first_name, last_name, username, lang_code, is_bot, created) VALUES (?, ?, ?, ?, ?, ?, ?)", usersTableName)
	log.Print("createUser user")
	log.Print(fmt.Sprintf("sqlQuery: %s", sqlQuery))
	result, err := r.db.Exec(sqlQuery, user.TgId, user.FirstName, user.LastName, user.Username, user.LanguageCode, user.IsBot, time.Now().UTC())
	if err != nil {
		log.Printf("error inserting user into table `%s` : %s", usersTableName, err.Error())
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("error getting inserted id %s", err.Error())
		return -1, err
	}
	return id, nil
}

func (r Repository) FindUserById(id int64) (User, error) {
	sqlQuery := fmt.Sprintf("SELECT id, tg_id, first_name, last_name, username, lang_code, is_bot FROM %s WHERE id = ?", usersTableName)
	rows := r.db.QueryRow(sqlQuery, id)
	var user User
	if err := rows.Scan(
		&user.Id,
		&user.TgId,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.LanguageCode,
		&user.IsBot,
		&user.Username,
		&user.Created); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("findUserById %d: no such user", id)
			return user, fmt.Errorf("findUserById %d: no such user", id)
		}
		return user, fmt.Errorf("findUserById %d: %v", id, err)
	}
	return user, nil
}

func (r Repository) FindUserByTgId(tgId int64) (User, error) {
	sqlQuery := fmt.Sprintf("SELECT id, tg_id, first_name, last_name, username, lang_code, is_bot FROM %s WHERE tg_id = ?", usersTableName)
	rows := r.db.QueryRow(sqlQuery, tgId)
	var user User
	if err := rows.Scan(
		&user.Id,
		&user.TgId,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.LanguageCode,
		&user.IsBot,
		&user.Username,
		&user.Created); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("findUserById %d: no such user", tgId)
			return user, err
		}
		return user, fmt.Errorf("findUserById %d: %v", tgId, err)
	}
	return user, nil
}

func InitDb() Repository {
	db := connectDb()
	return Repository{db: db}
}
