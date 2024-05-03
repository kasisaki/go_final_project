package utils

import (
	"encoding/json"
	"net/http"
)

func WriteNormalResponse(w http.ResponseWriter, strct any) {
	responseJSON, err := json.Marshal(strct)
	if err != nil {
		HandleError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
