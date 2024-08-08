package main

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
	m  sync.Mutex
}

func NewStorage(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAllTasks() []Task {
	s.m.Lock()
	defer s.m.Unlock()

	rows, err := s.db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		fmt.Printf("Failed to get all tasks: %v\n", err)
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Done)
		if err != nil {
			fmt.Printf("Failed to scan row: %v\n", err)
			return nil
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Row iteration error: %v\n", err)
		return nil
	}

	return tasks
}

func (s *Storage) CreateOneTask(t Task) int {
	s.m.Lock()
	defer s.m.Unlock()

	var id int
	err := s.db.QueryRow("INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id", t.Title, t.Done).Scan(&id)
	if err != nil {
		fmt.Printf("Failed to insert task: %v\n", err)
		return 0
	}
	t.ID = id

	fmt.Printf("Created task. Last ID: %v\n", t.ID)
	return t.ID
}

func (s *Storage) UpdateTask(t Task) bool {
	s.m.Lock()
	defer s.m.Unlock()

	res, err := s.db.Exec("UPDATE tasks SET title = $1, done = $2 WHERE id = $3", t.Title, t.Done, t.ID)
	if err != nil {
		fmt.Printf("Failed to update task: %v\n", err)
		return false
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Failed to get rows affected: %v\n", err)
		return false
	}

	return rowsAffected > 0
}

func (s *Storage) GetTaskByID(id int) (Task, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	var t Task
	err := s.db.QueryRow("SELECT id, title, done FROM tasks WHERE id = $1", id).Scan(&t.ID, &t.Title, &t.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return Task{}, false
		}
		fmt.Printf("Failed to get task by ID: %v\n", err)
		return Task{}, false
	}

	return t, true
}

func (s *Storage) DeleteTaskByID(id int) bool {
	s.m.Lock()
	defer s.m.Unlock()

	res, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		fmt.Printf("Failed to delete task: %v\n", err)
		return false
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Failed to get rows affected: %v\n", err)
		return false
	}

	return rowsAffected > 0
}
