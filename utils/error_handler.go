package utils

import (
	"encoding/json"
	"go_final_project/models"
	"net/http"
)

func HandleError(w http.ResponseWriter, status int, err error) {
	errResponse := models.ErrorResponse{Error: err.Error()}
	responseJSON, marshalErr := json.Marshal(errResponse)
	if marshalErr != nil {
		http.Error(w, marshalErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(responseJSON)
}
