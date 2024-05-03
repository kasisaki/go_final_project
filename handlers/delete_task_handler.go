package handlers

import (
	"errors"
	"go_final_project/db"
	"go_final_project/utils"
	"net/http"
	"strconv"
)

func HandleDeleteTask(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, errors.New("Неверный идентификатор"))
		return
	}
	_, err = db.GetTaskById(id)
	if err != nil {
		utils.HandleError(w, http.StatusBadRequest, err)
		return
	}
	err = db.DeleteById(id)
	if err != nil {
		utils.HandleError(w, http.StatusInternalServerError, err)
		return
	}

	// Если задача успешно удалена
	utils.WriteNormalResponse(w, "")
	return
}
