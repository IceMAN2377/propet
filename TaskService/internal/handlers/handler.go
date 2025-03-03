package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func InitDB(connStr string) error {
	var err error
	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		return err
	}
	return db.Ping(context.Background())
}

type Task struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func (t *Task) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO tasks (title, status) VALUES ($1, $2) RETURNING id`
	err := db.QueryRow(r.Context(), query, task.Title, task.Status).Scan(&task.Id)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (t *Task) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	var task Task
	query := `SELECT id, title, status FROM tasks WHERE id = $1`
	err = db.QueryRow(r.Context(), query, id).Scan(&task.Id, &task.Title, &task.Status)
	if err != nil {
		log.Printf("Failed to get task: %v", err)
		http.Error(w, "Failed to get task", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (t *Task) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `UPDATE tasks SET title = $1, status = $2 WHERE id = $3`
	existCheck, err := db.Exec(r.Context(), query, task.Title, task.Status, id)
	if err != nil {
		log.Printf("Failed to update task: %v", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}
	if existCheck.RowsAffected() == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func (t *Task) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM tasks WHERE id = $1`
	existCheck, err := db.Exec(r.Context(), query, id)
	if err != nil {
		log.Printf("Failed to delete task: %v", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	if existCheck.RowsAffected() == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Task deleted successfully"})
}
