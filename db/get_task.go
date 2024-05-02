package db

import "go_final_project/models"

func GetTasksFromDB() ([]models.TaskDTO, error) {
	// Здесь ваш код для получения списка задач из базы данных
	// Пример кода:
	tasks := make([]models.TaskDTO, 0)
	getTasks := `
	SELECT id, title, date, comment, repeat FROM scheduler
	ORDER BY date
	`
	res, err := db.Query(getTasks)
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
