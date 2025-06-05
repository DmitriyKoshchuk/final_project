package api

import (
	"encoding/json"
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.HandleFunc("/api/task/done", doneTaskHandler)
}

func writeJson(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, statusCode int, err error) {
	writeJson(w, statusCode, map[string]string{
		"error": err.Error(),
	})
}
