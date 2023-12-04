package db

import (
	"apitest/models"
	"database/sql"
	"log"
)

func InitDb(db *sql.DB) error {
	statement, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT,
		password TEXT
	);
	`)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	users := []models.NewUser{{Login: "user1", Password: "secret"}, {Login: "user2", Password: "secret"}}
	for _, user := range users {
		log.Printf("inserting %v\n", user)
		statement, err := db.Prepare("INSERT INTO users (login, password) VALUES (?, ?)")
		if err != nil {
			log.Print(err)
			return err
		}
		_, err = statement.Exec(user.Login, user.Password)
		if err != nil {
			log.Print(err)
			return err
		}
	}

	return nil
}
