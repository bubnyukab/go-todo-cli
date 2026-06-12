package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

func New() *Store { return &Store{} }

func (s *Store) Init() error {
	dbPath, err := GetDatabasePath()
	if err != nil {
		log.Fatalf("unable to resolve database path: %v", err)
	}

	s.conn, err = sql.Open("sqlite3", dbPath)
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

// GetDatabasePath dynamically resolves the cross-platform path for todos.db
func GetDatabasePath() (string, error) {
	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not detect user config directory: %w", err)
	}

	appDir := filepath.Join(baseDir, "todo")

	err = os.MkdirAll(appDir, 0o755)
	if err != nil {
		return "", fmt.Errorf("could not create application directory: %w", err)
	}

	return filepath.Join(appDir, "todos.db"), nil
}
