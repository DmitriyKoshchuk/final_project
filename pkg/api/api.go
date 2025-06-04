package api

import (
	"encoding/json"
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("/api/task", taskHandler)
	mux.HandleFunc("/api/nextdate", nextDateHandler)
	mux.HandleFunc("/api/tasks", tasksHandler)
	mux.HandleFunc("/api/task/done", doneTaskHandler) // новый обработчик
}

func writeJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	writeJson(w, map[string]string{
		"error": err.Error(),
	})
}
