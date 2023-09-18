package repository

import (
	"database/sql"
	"errors"

	"github.com/Palolem144/todo_list/internal/domain"
	"github.com/Palolem144/todo_list/internal/storage"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(db *sql.DB, task domain.Task) (*domain.Task, error) {
	if task.Name == "" {
		return &domain.Task{}, ErrInvalidTask
	}
	rows, err := storage.DB.Exec("INSERT INTO tasks(name) values (?)", task.Name)
	if err != nil {
		return &domain.Task{}, err
	}

	task.Id, err = rows.LastInsertId()
	if err != nil {
		return &domain.Task{}, err
	}
	return &task, nil
}

func (r *Repository) Get(db *sql.DB, id int64) (*domain.Task, error) {
	row := storage.DB.QueryRow("SELECT id, name from tasks where id=$1", id)
	var task domain.Task

	if err := row.Scan(&task.Id, &task.Name); err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *Repository) GetAll(db *sql.DB, id int64, name string) ([]domain.Task, error) {
	rows, err := storage.DB.Query("SELECT id, name FROM tasks")
	if err != nil {
		return []domain.Task{}, err
	}
	var tasks []domain.Task
	defer rows.Close()

	for rows.Next() {
		var task domain.Task
		err = rows.Scan(&task.Id, &task.Name)
		if err != nil {
			return []domain.Task{}, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *Repository) Update(db *sql.DB, task domain.Task) (*domain.Task, error) {
	rows, err := storage.DB.Exec("Update tasks SET name = ? WHERE id = ?", task.Name, task.Id)
	if err != nil {
		return &domain.Task{}, err
	}

	rowAffected, err := rows.RowsAffected()
	if err != nil {
		return &domain.Task{}, err
	}
	if rowAffected == 0 {
		return &domain.Task{}, ErrNoFoundTask
	}
	return &task, nil
}

func (r *Repository) Delete(db *sql.DB, id int64) (int64, error) {
	stmt, err := storage.DB.Exec("DELETE FROM tasks where id = ?", id)
	if err != nil {
		return 0, err
	}
	rowAffected, err := stmt.RowsAffected()
	if err != nil {
		return 0, err
	}
	if rowAffected == 0 {
		return 0, ErrNoFoundTask
	}
	return rowAffected, nil
}

var ErrNoFoundTask = errors.New("task not found")
var ErrInvalidTask = errors.New("TaskName error: name field  should not be empty")
