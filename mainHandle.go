package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

func handleTime(res http.ResponseWriter, req *http.Request) {
	s := time.Now().Format("02.01.2006 15:04:05")
	res.Write([]byte(s))
}

func handleMain(res http.ResponseWriter, req *http.Request) {
	webDir := "./web/"
	readFile, _ := os.ReadFile(webDir)

	res.Write(readFile)
}

func main() {
	setupDb()

	webDir := "web"
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Println("No PORT number provided")
		port = "7540"
	}
	db, exists := os.LookupEnv("TODO_DBFILE")
	fmt.Println(db)
	if !exists {
		log.Println("No BD path provided")
		db = "scheduler.db"
	}
	r := chi.NewRouter()

	r.Get("/time", handleTime)
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
