package repository

import (
	"database/sql"
	"github.com/Palolem144/todo_list/internal/domain"
)

// type Taskist interface {
// 	Create(userId int, tasks []domain.Task) (int, error)
// }

type Repository struct{}

func (r *Repository) GetByID(db *sql.DB, id int64) (*domain.Task, error) {
	row := db.QueryRow("SELECT * FROM tasks where id = ?", id)

	var task domain.Task
	if err := row.Scan(&task.Id, &task.Name); err != nil {
		return nil, err
	}

	return &task, nil
}
