package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var MainDB *sql.DB

// const FileName = "todo.db"

func InitDB() {
	// db_file = ".todo_list/todo.db"
	db, err := sql.Open("sqlite3", FileName)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	log.Println("Connected to the database")
	MainDB = db

}
