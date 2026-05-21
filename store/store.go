package store

import (
	"database/sql"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	ID   uuid.UUID
	Body string
	Done bool
}

func (t Todo) FilterValue() string { return t.Body }

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.conn, err = sql.Open("sqlite3", "./todos.db")
	if err != nil {
		return err
	}

	createTableStmt := `CREATE TABLE IF NOT EXISTS todos (
		id text not null primary key,
		body text not null,
		done boolean not null default 0
	);`

	if _, err := s.conn.Exec(createTableStmt); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetTodos() ([]Todo, error) {
	rows, err := s.conn.Query("SELECT id, body, done FROM todos")
	if err != nil {
		return nil, err
	}

	todos := []Todo{}
	defer rows.Close()
	for rows.Next() {
		todo := Todo{}
		rows.Scan(&todo.ID, &todo.Body, &todo.Done)
		todos = append(todos, todo)
	}

	return todos, nil
}

func (s *Store) SaveTodo(todo Todo) error {
	if todo.ID == uuid.Nil {
		todo.ID = uuid.New()
	}

	upsertQuery := `
	INSERT INTO todos (id, body, done)
	VALUES (?, ?, ?)
	ON CONFLICT(id) DO UPDATE
	SET body=excluded.body, done=excluded.done
	`

	if _, err := s.conn.Exec(upsertQuery, todo.ID, todo.Body, todo.Done); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteTodo(id uuid.UUID) error {
	_, err := s.conn.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}
