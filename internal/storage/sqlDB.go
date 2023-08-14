package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	db, err := sql.Open("sqlite3", "todo.db")
	if err != nil {
		log.Println("Cannot connect to db", err)
	}

	sqlDB := `
	CREATE TABLE IF NOT EXISTS tasks(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NULL
	);
	`
	_, err = db.Exec(sqlDB)
	if err != nil {
		log.Println("Cannot create table")
	}

	log.Println("Connected to the database")
	DB = db
}
