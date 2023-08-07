package models

import (
	"context"
	"time"
)

type ToDo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) GetToDo(id int) (ToDo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todo ToDo

	row := m.DB.QueryRowContext(ctx, "SELECT id, name, completed, created_at, updated_at FROM todos where id = $1", id)
	err := row.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return todo, err
	}

	return todo, err
}

func (m *DBModel) GetToDos() ([]ToDo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todos []ToDo

	rows, err := m.DB.QueryContext(ctx, "SELECT id, name, completed, created_at, updated_at FROM todos")

	if err != nil {
		return todos, err
	}

	defer rows.Close()

	for rows.Next() {
		var todo ToDo
		if err := rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return todos, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return todos, err
	}

	return todos, nil
}

func (m *DBModel) CreateToDo(t ToDo) (ToDo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into todos (name, completed, created_at, updated_at)
		values ($1, $2, $3, $4)
		returning id
	`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt, t.Name, t.Completed, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return ToDo{}, err
	}
	t.ID = id // Assuming ToDo has an ID field
	return t, nil
}

func (m *DBModel) EditToDo(t ToDo) (ToDo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todo ToDo

	updateStmt := `
		update todos set 
			name = $1,
			completed = $2,
			updated_at = $3
			where 
			id = $4
		`
	_, err := m.DB.ExecContext(ctx, updateStmt, t.Name, t.Completed, time.Now(), t.ID)

	if err != nil {
		return todo, err
	}
	selectStmt := `select id, name, completed, created_at, updated_at from todos where id = $1`
	err = m.DB.QueryRowContext(ctx, selectStmt, t.ID).Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return todo, err
	}

	return todo, nil
}

func (m *DBModel) DeleteToDo(id int) (ToDo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var todo ToDo

	selectStmt := `select id, name, completed, created_at, updated_at from todos where id = $1`
	err := m.DB.QueryRowContext(ctx, selectStmt, id).Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return todo, err
	}

	deleteStmt := `delete from todos where id = $1`
	_, err = m.DB.ExecContext(ctx, deleteStmt, id)

	if err != nil {
		return todo, err
	}

	return todo, nil
}
