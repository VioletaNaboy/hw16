package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func main() {
	mux := mux.NewRouter()

	connStr := os.Getenv("POSTGRES_CONN_STR")
	storage, err := NewStorage(connStr)
	if err != nil {
		log.Fatal().Msgf("Failed to create storage: %v", err)
	}

	tasks := TaskResource{
		s: storage,
	}

	mux.HandleFunc("/tasks", tasks.GetAll).Methods("GET")
	mux.HandleFunc("/tasks/add", tasks.CreateOne).Methods("POST")
	mux.HandleFunc("/tasks/update", tasks.UpdateOne).Methods("PUT")
	mux.HandleFunc("/tasks/delete/{id:[0-9]+}", tasks.DeleteOne).Methods("DELETE")

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
	vars := mux.Vars(r)
	idStr := vars["id"]
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
