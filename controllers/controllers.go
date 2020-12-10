package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/AthanatiusC/TaskManager/models"
)

func respondJSON(w http.ResponseWriter, status int, message string, data interface{}) {
	var payload models.Payload
	if status == 200 {
		payload.Status = true
	} else {
		payload.Status = false
	}
	payload.Message = message
	payload.Data = data

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
