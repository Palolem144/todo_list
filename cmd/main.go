package main

import (
	"log"
	"net/http"

	"github.com/Palolem144/todo_list/internal/handler"
	"github.com/Palolem144/todo_list/internal/models"
)

func main() {

	models.InitDB()
	http.HandleFunc("/create", handler.Create)
	http.HandleFunc("/getAll", handler.GetAll)
	http.HandleFunc("/get", handler.Get)
	http.HandleFunc("/update", handler.Update)
	http.HandleFunc("/delete", handler.Delete)

	log.Println("Starting server...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
