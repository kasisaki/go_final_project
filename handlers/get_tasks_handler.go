package handlers

import (
	"encoding/json"
	"errors"
	"go_final_project/constants"
	"go_final_project/db"
	"go_final_project/models"
	"go_final_project/utils"
	"net/http"
	"strconv"
)

func HandleGetTasks(w http.ResponseWriter, req *http.Request) {
	// Проверяем GET-запрос или нет
	if req.Method == http.MethodGet {
		search := req.URL.Query().Get("search")
		tasks, err := db.GetTasksFromDB(search, constants.TasksNumberLimit) // Получаем список задач из базы данных лим
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, err)
			return
		}
		respMap := make(map[string][]models.TaskDTO)
		respMap["tasks"] = tasks

		// Преобразуем список задач в JSON
		responseJSON, err := json.Marshal(respMap)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, err)
			return
		}

		// Отправляем JSON клиенту
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}

func HandleGetTaskById(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		id := req.URL.Query().Get("id")
		_, err := strconv.Atoi(id)
		if err != nil {
			utils.HandleError(w, http.StatusBadRequest, errors.New("Не указан идентификатор"))
			return
		}
		task, err := db.GetTaskById(id)
		if err != nil {
			utils.HandleError(w, 404, err)
		}

		responseJSON, err := json.Marshal(task)
		// Отправляем JSON клиенту
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}
