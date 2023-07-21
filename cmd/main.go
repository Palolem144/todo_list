package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Palolem144/todo_list/internal/handler"
	"github.com/Palolem144/todo_list/internal/storage"
)

func main() {
	defaultPort := ":8000"
	storage.InitDB()
	http.HandleFunc("/create", handler.Create)
	http.HandleFunc("/getAll", handler.GetAll)
	http.HandleFunc("/get", handler.Get)
	http.HandleFunc("/update", handler.Update)
	http.HandleFunc("/delete", handler.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Println("Starting server...")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
