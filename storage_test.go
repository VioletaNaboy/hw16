package main

import (
	"testing"
)

func TestStorage_CreateOneTask(t *testing.T) {
	storage := NewStorage()
	task := Task{Title: "Test Task", Done: false}

	id := storage.CreateOneTask(task)
	if id != 1 {
		t.Fatalf("Expected ID to be 1, got %d", id)
	}

	storedTask, exists := storage.allTasks[id]
	if !exists {
		t.Fatal("Expected task to be stored, but it was not found")
	}

	if storedTask.Title != task.Title || storedTask.Done != task.Done {
		t.Fatalf("Stored task does not match original task. Got %+v, want %+v", storedTask, task)
	}
}

func TestStorage_UpdateTask(t *testing.T) {
	storage := NewStorage()
	task := Task{Title: "Test Task", Done: false}
	id := storage.CreateOneTask(task)

	updatedTask := Task{ID: id, Title: "Updated Task", Done: true}
	success := storage.UpdateTask(updatedTask)
	if !success {
		t.Fatal("Expected UpdateTask to succeed, but it failed")
	}

	storedTask, exists := storage.allTasks[id]
	if !exists {
		t.Fatal("Expected task to be stored, but it was not found")
	}

	if storedTask.Title != updatedTask.Title || storedTask.Done != updatedTask.Done {
		t.Fatalf("Updated task does not match expected task. Got %+v, want %+v", storedTask, updatedTask)
	}
}

func TestStorage_DeleteTaskByID(t *testing.T) {
	storage := NewStorage()
	task := Task{Title: "Test Task", Done: false}
	id := storage.CreateOneTask(task)

	success := storage.DeleteTaskByID(id)
	if !success {
		t.Fatal("Expected DeleteTaskByID to succeed, but it failed")
	}

	_, exists := storage.allTasks[id]
	if exists {
		t.Fatal("Expected task to be deleted, but it was still found")
	}
}

func TestStorage_UpdateTask_FailOnNonExistentTask(t *testing.T) {
	s := NewStorage()

	task := Task{ID: 999, Title: "Updated Task", Done: true}
	success := s.UpdateTask(task)

	if success {
		t.Errorf("Expected update to fail for non-existent task")
	}
}

func TestStorage_DeleteTaskByID_FailOnNonExistentTask(t *testing.T) {
	s := NewStorage()

	success := s.DeleteTaskByID(999)

	if success {
		t.Errorf("Expected deletion to fail for non-existent task")
	}
}
