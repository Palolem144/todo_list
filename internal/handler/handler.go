package handler

import (
	"encoding/json"
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
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer rows.Close()
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		tasks = append(tasks, task)
	}
	respBody, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
	ID := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	stmt := storage.DB.QueryRow("SELECT * FROM tasks where id = ?", id)
	var task domain.Task

	err = stmt.Scan(&task.Id, &task.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	respBody, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := r.FormValue("name")
	var task domain.Task
	task.Name = name

	stmt, err := storage.DB.Exec("INSERT INTO tasks(name) values (?)", task.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	newID, err := stmt.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	task.Id = newID
	resBody, err := json.Marshal(task)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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
	ID := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	var task domain.Task
	task.Id = id
	task.Name = name

	stmt, err := storage.DB.Exec("UPDATE tasks SET name = ? WHERE id = ?", task.Name, task.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	rowAffected, err := stmt.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
	respBody, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(respBody)

}

func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	ID := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	stmt, err := storage.DB.Exec("DELETE FROM tasks where id = ?", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	rowAffected, err := stmt.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	msg := strconv.Itoa(int(rowAffected))
	w.Write([]byte(msg))
}
