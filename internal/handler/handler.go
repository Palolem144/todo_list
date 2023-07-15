package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Palolem144/todo_list/internal/domain"
	"github.com/Palolem144/todo_list/internal/models"
)

type Tasks []domain.Task

func GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	rows, err := models.DB.Query("SELECT * FROM tasks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var tasks Tasks
	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Name)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	respBody, errMarshal := json.Marshal(tasks)
	if errMarshal != nil {
		log.Fatal(errMarshal)
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
	id := r.URL.Query().Get("id")
	stmt, err := models.DB.Prepare(" SELECT * FROM tasks where id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	rows, errQuery := stmt.Query(id)
	if errQuery != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	var task domain.Task
	for rows.Next() {
		err = rows.Scan(&task.Id, &task.Name)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		respBody, errMarshal := json.Marshal(task)
		if errMarshal != nil {
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
	stmt, err := models.DB.Prepare("INSERT INTO tasks(name) values (?)")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	result, errExec := stmt.Exec(task.Name)
	if errExec != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	newID, errLast := result.LastInsertId()
	if errLast != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	task.Id = newID
	resBody, errMarshal := json.Marshal(task)

	if errMarshal != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = json.NewDecoder(r.Body).Decode(task.Name)
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
	ID, _ := strconv.ParseInt(id, 10, 0)
	task.Id = ID
	task.Name = name
	stmt, err := models.DB.Prepare("UPDATE  tasks SET name = ? WHERE id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	result, errExec := stmt.Exec(task.Name, task.Id)
	if errExec != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	rowAffected, errAff := result.RowsAffected()
	if errAff != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if rowAffected > 0 {
		respBody, errMarshal := json.Marshal(task)
		if errMarshal != nil {
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
	stmt, err := models.DB.Prepare("DELETE FROM tasks where id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	result, errExec := stmt.Exec(id)
	if errExec != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	rowAffected, errRow := result.RowsAffected()
	if errRow != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	msg := strconv.Itoa(int(rowAffected))
	w.Write([]byte(msg))
}
