package db

import (
	"errors"
	"go_final_project/models"
	"strconv"
)

func PutTask(task models.TaskDTO) error {
	// Проверяем, существует ли задача с указанным ID
	var id int
	if err := DB.QueryRow("SELECT id FROM scheduler WHERE id = ?", task.Id).Scan(&id); err != nil {
		return errors.New("Задача не найдена")
	}

	if id == 0 {
		return errors.New("Задача не найдена")
	}

	// Выполняем запрос UPDATE в рамках одной транзакции
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			return
		}
		err = tx.Commit()
	}()
	query := `
        UPDATE scheduler
        SET date = ?, title = ?, comment = ?, repeat = ? 
        WHERE id = ?
    `

	idInt, _ := strconv.Atoi(task.Id)

	_, err = tx.Exec(query,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		idInt,
	)
	if err != nil {
		return err
	}
	return nil
}
