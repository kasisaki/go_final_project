package models

import "fmt"

type ErrorResponse struct {
	Error string `json:"error"`
}

// TaskDTO описывает структуры задачи
type TaskDTO struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Task struct {
	Id      int64  `db:"id,omitempty"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

func (t TaskDTO) toString() string {
	return fmt.Sprintf("ID: %d,\nTitle %s,\nDate: %s,\nComment: %s,\nRepeat: %s",
		t.Id, t.Title, t.Date, t.Comment, t.Repeat)
}
