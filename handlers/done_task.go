package handlers

import (
	"errors"
	"go_final_project/db"
	"go_final_project/models"
	"go_final_project/services"
	"net/http"
	"strconv"
	"time"
)

func HandleTaskDone(w http.ResponseWriter, req *http.Request) {
	var task models.TaskDTO
	id, err := strconv.Atoi(req.URL.Query().Get("id"))

	if err != nil {
		HandleError(w, http.StatusBadRequest, errors.New("Неверный идентификатор"))
		return
	}
	task, err = db.GetTaskById(id)
	if err != nil {
		HandleError(w, http.StatusBadRequest, err)
		return
	}
	if task.Repeat == "" {
		err = db.DeleteById(id)
		if err != nil {
			HandleError(w, http.StatusInternalServerError, err)
			return
		}
		HandleNormalResponse(w, "")
		return
	}
	now := time.Now().Truncate(24 * time.Hour)
	task.Date, err = services.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err)
		return
	}
	err = db.PutTask(task)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err)
		return
	}

	HandleNormalResponse(w, "")
	return
}

func HandleDeleteTask(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil {
		HandleError(w, http.StatusBadRequest, errors.New("Неверный идентификатор"))
		return
	}

	// Данная проверка проводится, чтобы вернуть ошибку, при попытке удаления несуществующей задачи
	_, err = db.GetTaskById(id)
	if HandleGetError(w, err) {
		return
	}

	err = db.DeleteById(id)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err)
		return
	}

	// Если задача успешно удалена
	HandleNormalResponse(w, "")
	return
}
