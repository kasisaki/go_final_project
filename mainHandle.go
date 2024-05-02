package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

var db *sql.DB

// Task описывает структуры задачи
type Task struct {
	Id      int64  `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (t Task) toString() string {
	return fmt.Sprintf("ID: %d,\nTitle %s,\nDate: %s,\nComment: %s,\nRepeat: %s",
		t.Id, t.Title, t.Date, t.Comment, t.Repeat)
}

func handleNextDate(res http.ResponseWriter, req *http.Request) {
	nowDate, _ := time.Parse(DateLayout, req.URL.Query().Get("now"))
	taskDate := req.URL.Query().Get("date")
	repeatRule := req.URL.Query().Get("repeat")
	nextDate, _ := NextDate(nowDate, taskDate, repeatRule)
	fmt.Printf("Now %s, taskDate %s, repeat %s, nextDate %s\n", nowDate.Format(DateLayout), taskDate, repeatRule, nextDate)
	res.Write([]byte(nextDate))
}

func handleError(w http.ResponseWriter, status int, err error) {
	errResponse := ErrorResponse{Error: err.Error()}
	responseJSON, marshalErr := json.Marshal(errResponse)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(responseJSON)
}

func handleTaskPost(w http.ResponseWriter, req *http.Request) {
	// Проверяем POST-запрос или нет
	if req.Method == http.MethodPost {
		var task Task
		var buf bytes.Buffer

		// читаем тело запроса
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// десериализуем JSON в Task
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(task.Title) == "" {
			// Если title пустой, возвращаем ошибку
			handleError(w, http.StatusBadRequest, errors.New("Не указан заголовок задачи"))
			return
		}

		now := time.Now().Truncate(24 * time.Hour)
		if task.Date == "" {
			task.Date = now.Format(DateLayout)
		}
		taskDate, err := time.Parse(DateLayout, task.Date)
		if err != nil {
			handleError(w, http.StatusBadRequest, err)
			return
		}

		if taskDate.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(DateLayout)
			} else {
				nextDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					handleError(w, http.StatusBadRequest, err)
					return
				}
				task.Date = nextDate
			}
		}

		id, err := insertTask(task)
		if err != nil {
			// Если произошла ошибка при вставке задачи
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		// Если задача успешно создана
		response := struct {
			ID int64 `json:"id"`
		}{ID: id}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			handleError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func insertTask(task Task) (int64, error) {
	insertSQL := `
	INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)
	`
	res, err := db.Exec(insertSQL, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Printf("Inserted task with ID: %d\n", lastInsertID)
	return lastInsertID, nil
}

func main() {
	db = setupDb()

	webDir := "web"
	port, exists := os.LookupEnv("PORT")
	if !exists {
		log.Println("No PORT number provided... Setting to default")
		port = "7540"
	}

	r := chi.NewRouter()

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))
	r.Get("/api/nextdate", handleNextDate)
	r.Post("/api/task", handleTaskPost)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}

// to remove

func checkDBConnection() error {
	if db == nil {
		return errors.New("database connection is not initialized")
	}
	fmt.Println("Trying to PING DB")

	err := db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}

	return nil
}
