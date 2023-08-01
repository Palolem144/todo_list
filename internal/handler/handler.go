package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Palolem144/todo_list/internal/domain"
	"github.com/Palolem144/todo_list/internal/storage"
)

func GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rows, err := storage.DB.Query("SELECT id, name FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Name)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	respBody, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := storage.DB.Prepare(" SELECT * FROM tasks where id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	rows, err := stmt.Query(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	defer rows.Close()
	var task domain.Task
	for rows.Next() {
		err = rows.Scan(&task.Id, &task.Name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		respBody, err := json.Marshal(task)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("name")
	var task domain.Task
	task.Name = name
	stmt, err := storage.DB.Prepare("INSERT INTO tasks(name) values (?)")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer stmt.Close()
	result, err := stmt.Exec(task.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	newID, err := result.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	task.Id = newID
	resBody, err := json.Marshal(task)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("name")
	id := r.URL.Query().Get("id")
	var task domain.Task
	ID, _ := strconv.ParseInt(id, 10, 64)
	task.Id = ID
	task.Name = name
	stmt, err := storage.DB.Prepare("UPDATE  tasks SET name = ? WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer stmt.Close()
	result, err := stmt.Exec(task.Name, task.Id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if rowAffected > 0 {
		respBody, err := json.Marshal(task)
		if err != nil {
			w.WriteHeader(http.StatusNotModified)
		}
		w.Write(respBody)
	} else {
		msg := strconv.Itoa(int(rowAffected))
		w.Write([]byte(msg))
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	id := r.URL.Query().Get("id")
	stmt, err := storage.DB.Prepare("DELETE FROM tasks where id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	defer stmt.Close()
	result, err := stmt.Exec(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	msg := strconv.Itoa(int(rowAffected))
	w.Write([]byte(msg))
}
