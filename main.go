package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	mux := http.NewServeMux()

	s := NewStorage()

	tasks := TaskResource{
		s: s,
	}

	mux.HandleFunc("GET /tasks", tasks.GetAll)
	mux.HandleFunc("POST /tasks/add", tasks.CreateOne)

	//обовʼязково передати id
	mux.HandleFunc("PUT /tasks/update", tasks.UpdateOne)
	mux.HandleFunc("DELETE /tasks/delete", tasks.DeleteOne)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Failed to listen and serve: %v\n", err)
	}
}

type TaskResource struct {
	s *Storage
}

func (t *TaskResource) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks := t.s.GetAllTasks()

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *TaskResource) CreateOne(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Request body is missing", http.StatusBadRequest)
		return
	}
	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task.ID = t.s.CreateOneTask(task)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *TaskResource) UpdateOne(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Request body is missing", http.StatusBadRequest)
		return
	}
	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		fmt.Printf("Failed to decode: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	success := t.s.UpdateTask(task)
	if !success {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *TaskResource) DeleteOne(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	success := t.s.DeleteTaskByID(id)
	if !success {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
