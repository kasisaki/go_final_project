package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_final_project/constants"
	"go_final_project/db"
	"go_final_project/models"
	"go_final_project/services"
	"go_final_project/utils"
	"net/http"
	"strings"
	"time"
)

func HandleTaskPost(w http.ResponseWriter, req *http.Request) {
	// Проверяем POST-запрос или нет
	if req.Method == http.MethodPost {
		var task models.TaskDTO
		var buf bytes.Buffer

		// читаем тело запроса
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// десериализуем JSON в TaskDTO
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(task.Title) == "" {
			// Если title пустой, возвращаем ошибку
			utils.HandleError(w, http.StatusBadRequest, errors.New("Не указан заголовок задачи"))
			return
		}

		now := time.Now().Truncate(24 * time.Hour)
		if task.Date == "" {
			task.Date = now.Format(constants.DateLayout)
		}
		taskDate, err := time.Parse(constants.DateLayout, task.Date)
		if err != nil {
			utils.HandleError(w, http.StatusBadRequest, err)
			return
		}

		if taskDate.Before(now) {
			if task.Repeat == "" {
				task.Date = now.Format(constants.DateLayout)
			} else {
				nextDate, err := services.NextDate(now, task.Date, task.Repeat)
				if err != nil {
					utils.HandleError(w, http.StatusBadRequest, err)
					return
				}
				task.Date = nextDate
			}
		}

		id, err := db.InsertTask(task)
		if err != nil {
			// Если произошла ошибка при вставке задачи
			utils.HandleError(w, http.StatusInternalServerError, err)
			return
		}

		// Если задача успешно создана
		response := struct {
			ID int64 `json:"id"`
		}{ID: id}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			utils.HandleError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)
	}
}
