package handlers

import (
	"errors"
	"go_final_project/db"
	"go_final_project/models"
	"go_final_project/services"
	"go_final_project/utils"
	"net/http"
	"strconv"
	"time"
)

func HandleTaskDone(w http.ResponseWriter, req *http.Request) {
	var task models.TaskDTO
	id, err := strconv.Atoi(req.URL.Query().Get("id"))

	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, errors.New("Неверный идентификатор"))
		return
	}
	task, err = db.GetTaskById(id)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err)
		return
	}
	if task.Repeat == "" {
		err = db.DeleteById(id)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, err)
			return
		}
		utils.WriteNormalResponse(w, "")
		return
	}
	now := time.Now().Truncate(24 * time.Hour)
	task.Date, err = services.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err)
		return
	}
	err = db.PutTask(task)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteNormalResponse(w, "")
	return
}
