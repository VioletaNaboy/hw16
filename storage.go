package main

import (
	"fmt"
	"sort"
	"sync"
)

type Storage struct {
	m        sync.Mutex
	lastID   int
	allTasks map[int]Task
}

func NewStorage() *Storage {
	return &Storage{
		allTasks: make(map[int]Task),
	}
}

func (s *Storage) GetAllTasks() []Task {
	s.m.Lock()
	defer s.m.Unlock()

	var tasks = make([]Task, 0, len(s.allTasks))

	for _, t := range s.allTasks {
		tasks = append(tasks, t)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks
}

func (s *Storage) CreateOneTask(t Task) int {
	s.m.Lock()
	defer s.m.Unlock()

	fmt.Println("Trying to create task")
	t.ID = s.lastID + 1
	s.allTasks[t.ID] = t
	s.lastID++
	fmt.Printf("Created task. Last ID: %v\n", s.lastID)
	return t.ID
}

func (s *Storage) UpdateTask(t Task) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, ok := s.allTasks[t.ID]
	if !ok {
		return false
	}

	s.allTasks[t.ID] = t
	return true
}

func (s *Storage) GetTaskByID(id int) (Task, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	task, ok := s.allTasks[id]
	return task, ok
}

func (s *Storage) DeleteTaskByID(id int) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, ok := s.allTasks[id]
	if !ok {
		return false
	}

	delete(s.allTasks, id)
	return true
}
