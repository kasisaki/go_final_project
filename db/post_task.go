package db

import (
	"fmt"
	"go_final_project/models"
)

func InsertTask(task models.TaskDTO) (int64, error) {
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
