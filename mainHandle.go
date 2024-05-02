package main

import (
	"github.com/go-chi/chi/v5"
	"go_final_project/constants"
	"go_final_project/db"
	"go_final_project/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	db.SetupDb()

	webDir := "web"
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Println("No PORT number provided... Setting to default")
		port = constants.DefaultPort
	}

	r := chi.NewRouter()

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
	r.Get("/api/nextdate", handlers.HandleNextDate)
	r.Post("/api/task", handlers.HandleTaskPost)
	r.Get("/api/tasks", handlers.HandleGetTasks)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
