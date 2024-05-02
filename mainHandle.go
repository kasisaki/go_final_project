package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

// Task описывает структуры задачи
type Task struct {
	ID      int       `json:"id"`
	Date    time.Time `json:"date"`
	Title   string    `json:"title"`
	Comment string    `json:"comment"`
	Repeat  string    `json:"repeat"`
}

func handleNextDate(res http.ResponseWriter, req *http.Request) {
	nowDate, _ := time.Parse(DateLayout, req.URL.Query().Get("now"))
	taskDate := req.URL.Query().Get("date")
	repeatRule := req.URL.Query().Get("repeat")
	nextDate, _ := NextDate(nowDate, taskDate, repeatRule)
	fmt.Printf("Now %s, taskDate %s, repeat %s, nextDate %s\n", nowDate.Format(DateLayout), taskDate, repeatRule, nextDate)
	res.Write([]byte(nextDate))
}

func main() {
	setupDb()

	webDir := "web"
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Println("No PORT number provided... Setting to default")
		port = "7540"
	}
	dbName, exists := os.LookupEnv("TODO_DBFILE")
	fmt.Println(dbName)
	if !exists {
		log.Println("No BD path provided... Setting to default")
		dbName = "scheduler.db"
	}
	r := chi.NewRouter()

	r.Get("/api/nextdate", handleNextDate)
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
