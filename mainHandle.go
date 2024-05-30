package main

import (
	"github.com/go-chi/chi/v5"
	"go_final_project/db"
	"go_final_project/handlers"
	"log"
	"net/http"
	"os"
)

const DefaultPort = "7540"

func main() {
	db.SetupDb()

	webDir := "web"
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Println("No PORT number provided... Setting to default")
		port = DefaultPort
	}

	r := chi.NewRouter()

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
	r.Get("/api/nextdate", handlers.HandleNextDate)
	r.Get("/api/tasks", handlers.HandleGetTasks)
	r.Get("/api/task", handlers.HandleGetTaskById)
	r.Post("/api/task", handlers.HandlePostTask)
	r.Post("/api/task/done", handlers.HandleTaskDone)
	r.Put("/api/task", handlers.HandlePutTask)
	r.Delete("/api/task", handlers.HandleDeleteTask)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
