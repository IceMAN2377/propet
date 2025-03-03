package main

import (
	"log"
	"net/http"

	"github.com/IceMAN2377/appl/internal/handlers"
)

func main() {
	connStr := "postgres://postgres:mypassword@localhost:5432/psdb?sslmode=disable"
	if err := handlers.InitDB(connStr); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	router := http.NewServeMux()
	taskHandler := handlers.Task{}
	router.HandleFunc("POST /tasks", taskHandler.CreateTask)
	router.HandleFunc("GET /tasks/{id}", taskHandler.GetTask)
	router.HandleFunc("PUT /tasks/{id}", taskHandler.UpdateTask)
	router.HandleFunc("DELETE /tasks/{id}", taskHandler.DeleteTask)

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
