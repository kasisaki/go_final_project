package db

import (
	"database/sql"
	"errors"
	"go_final_project/models"
	"go_final_project/utils"
)

func GetTasksFromDB(search string, limit int) ([]models.TaskDTO, error) {
	// Здесь ваш код для получения списка задач из базы данных
	// Пример кода:
	tasks := make([]models.TaskDTO, 0)
	isDate, date := utils.ValidateDate(search)
	var getTasks string
	if search == "" {
		getTasks = `
	SELECT id, title, date, comment, repeat FROM scheduler
	ORDER BY date
	LIMIT :limit
	`
	} else if isDate {
		getTasks = `
	SELECT id, title, date, comment, repeat FROM scheduler
	WHERE date = :date
	ORDER BY date
	LIMIT :limit
	`
	} else {
		getTasks = `
	SELECT id, title, date, comment, repeat FROM scheduler
	WHERE title LIKE '%' || :search || '%' OR comment LIKE '%' || :search || '%'
	ORDER BY date
	LIMIT :limit
	`
	}
	res, err := DB.Query(getTasks,
		sql.Named("search", search), sql.Named("limit", limit), sql.Named("date", date))
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var task models.TaskDTO
		if err := res.Scan(&task.Id, &task.Title, &task.Date, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func GetTaskById(id int) (models.TaskDTO, error) {
	var task models.TaskDTO

	tx, err := DB.Begin()
	if err != nil {
		return task, err
	}
	defer tx.Rollback()

	query := `SELECT id, title, date, comment, repeat FROM scheduler WHERE id = ?`
	res, err := tx.Query(query, id)
	if err != nil {
		return task, err
	}
	defer res.Close()

	if !res.Next() {
		return task, errors.New("Задача не найдена")
	}

	if err := res.Scan(&task.Id, &task.Title, &task.Date, &task.Comment, &task.Repeat); err != nil {
		return task, err
	}

	err = tx.Commit()
	if err != nil {
		return task, err
	}

	return task, nil
}
