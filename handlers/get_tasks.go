package handlers

import (
	"encoding/json"
	"go_final_project/db"
	"go_final_project/models"
	"go_final_project/utils"
	"net/http"
)

func HandleGetTasks(w http.ResponseWriter, req *http.Request) {
	// Проверяем GET-запрос или нет
	if req.Method == http.MethodGet {
		tasks, err := db.GetTasksFromDB() // Получаем список задач из базы данных
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
