package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Palolem144/todo_list/internal/domain"
	"github.com/Palolem144/todo_list/internal/repository"
)

type Repository interface {
	Create(db *sql.DB, Task domain.Task) (*domain.Task, error)
	Get(db *sql.DB, id int64) (*domain.Task, error)
	GetAll(db *sql.DB, id int64, name string) ([]domain.Task, error)
	Delete(db *sql.DB, id int64) (int64, error)
	Update(db *sql.DB, task domain.Task) (*domain.Task, error)
}
type TaskHandler struct {
	db *sql.DB
	rp Repository
}

func NewHandler(rp Repository) TaskHandler {
	return TaskHandler{
		rp: rp,
	}
}

func (th *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var task domain.Task
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&task)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	createdTask, err := th.rp.Create(th.db, task)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidTask) {
			log.Println(repository.ErrInvalidTask)
			errMsg := fmt.Sprintf("task name is empty: %s", task.Name)
			w.Write([]byte(errMsg))
			return
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	msg := fmt.Sprintf("Task id: %d, name: %s, created", createdTask.Id, createdTask.Name)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(msg))
}

func (th *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ID := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		log.Println("Parse id error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	task, err := th.rp.Get(th.db, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("Get task error", err)
		if errors.Is(err, sql.ErrNoRows) {
			errMsg := fmt.Sprintf("Task with id %d does not exist", id)
			w.Write([]byte(errMsg))
			return
		}
		return
	}
	respBody, err := json.Marshal(&task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}

func (th *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var task domain.Task
	defer r.Body.Close()

	tasks, err := th.rp.GetAll(th.db, task.Id, task.Name)
	if err != nil {
		log.Println("GetAllTasks error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respBody, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func (th *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	defer r.Body.Close()

	json.NewDecoder(r.Body).Decode(&task)

	_, err = th.rp.Update(th.db, task)
	if err != nil {
		if errors.Is(err, repository.ErrNoFoundTask) {
			log.Println(repository.ErrNoFoundTask)
			errMsg := fmt.Sprintf("task with id: %d, not found", task.Id)
			w.Write([]byte(errMsg))
			return
		}
		return
	}

	respBody, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)

}

func (th *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	ID := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	stmt, err := th.rp.Delete(th.db, id)
	if err != nil {
		if errors.Is(err, repository.ErrNoFoundTask) {
			log.Println(repository.ErrNoFoundTask)
			errMsg := fmt.Sprintf("task with id: %d, not found", id)
			w.Write([]byte(errMsg))
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	msg := fmt.Sprintf("%d task with id: %d, has been deleted", stmt, id)
	w.Write([]byte(msg))
}
