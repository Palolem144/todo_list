package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Palolem144/todo_list/internal/handler"
	"github.com/Palolem144/todo_list/internal/storage"
)

func main() {
	defaultAddr := ":8000"
	storage.InitDB()
	http.HandleFunc("/create", handler.Create)
	http.HandleFunc("/getAll", handler.GetAll)
	http.HandleFunc("/get", handler.Get)
	http.HandleFunc("/update", handler.Update)
	http.HandleFunc("/delete", handler.Delete)

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
