package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const FileName = "/Users/pllm/go/src/github.com/Palolem144/todo_list/todo.db"

func InitDB() {
	db, err := sql.Open("sqlite3", FileName)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()
	log.Println("Connected to the database")
	DB = db
}
