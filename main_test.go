package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestCreateOneTask(t *testing.T) {
	s := NewStorage()
	tasks := TaskResource{s: s}

	task := Task{Title: "Test Task", Done: false}
	body, _ := json.Marshal(task)
	req, err := http.NewRequest("POST", "/tasks/add", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tasks.CreateOne)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var createdTask Task
	if err := json.NewDecoder(rr.Body).Decode(&createdTask); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if createdTask.ID <= 0 {
		t.Errorf("Expected a valid ID, got %v", createdTask.ID)
	}

	if createdTask.Title != "Test Task" {
		t.Errorf("Expected Title to be 'Test Task', got %v", createdTask.Title)
	}
}

func TestUpdateTask(t *testing.T) {
	s := NewStorage()
	tasks := TaskResource{s: s}

	task := Task{Title: "Initial Task", Done: false}
	taskID := s.CreateOneTask(task)

	updatedTask := Task{ID: taskID, Title: "Updated Task", Done: true}
	body, _ := json.Marshal(updatedTask)
	req, err := http.NewRequest("PUT", "/tasks/update", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tasks.UpdateOne)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result Task
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.Title != "Updated Task" {
		t.Errorf("Expected Title to be 'Updated Task', got %v", result.Title)
	}

	if !result.Done {
		t.Errorf("Expected Done to be true, got false")
	}
}
func TestDeleteTask(t *testing.T) {
	s := NewStorage()
	tasks := TaskResource{s: s}

	task := Task{Title: "Task to Delete", Done: false}
	taskID := s.CreateOneTask(task)

	req, err := http.NewRequest("DELETE", "/tasks/delete?id="+strconv.Itoa(taskID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tasks.DeleteOne)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	if _, exists := s.GetTaskByID(taskID); exists {
		t.Errorf("Expected task to be deleted, but it still exists")
	}
}
func TestCreateOneTask_FailOnInvalidData(t *testing.T) {
	s := NewStorage()
	tasks := TaskResource{s: s}

	invalidTask := Task{}
	body, _ := json.Marshal(invalidTask)
	req, err := http.NewRequest("POST", "/tasks/add", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(tasks.CreateOne)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	var responseTask Task
	err = json.NewDecoder(rr.Body).Decode(&responseTask)
	if err == nil {
		t.Errorf("Expected error but got task with ID %v", responseTask.ID)
	}
}
