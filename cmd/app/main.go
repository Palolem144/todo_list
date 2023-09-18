package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Palolem144/todo_list/internal/handler"
	"github.com/Palolem144/todo_list/internal/repository"
	"github.com/Palolem144/todo_list/internal/storage"
)

func main() {

	dbPath := os.Getenv("DB_PATH")

	if dbPath == "" {
		dbPath = "todo.db"
	}
	handlers := handler.NewHandler(repository.NewRepository())

	defaultAddr := ":8080"
	storage.InitDB(dbPath)
	http.HandleFunc("/create", handlers.Create)
	http.HandleFunc("/getAll", handlers.GetAll)
	http.HandleFunc("/get", handlers.Get)
	http.HandleFunc("/update", handlers.Update)
	http.HandleFunc("/delete", handlers.Delete)

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = defaultAddr
	}

	log.Println("Starting server...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
