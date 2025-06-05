package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"go1f/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
	}
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJson(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Некорректные данные"})
		return
	}

	if task.ID == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if _, err := strconv.ParseInt(task.ID, 10, 64); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Неверный идентификатор"})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Пустой заголовок"})
		return
	}

	if !validDate(task.Date) {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Неверная дата"})
		return
	}

	if !validRepeat(task.Repeat) {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Неверный repeat"})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}

func validRepeat(s string) bool {
	if s == "" {
		return true
	}
	if s == "y" {
		return true
	}
	if len(s) > 2 && s[:2] == "d " {
		_, err := strconv.Atoi(s[2:])
		return err == nil
	}
	return false
}

func validDate(date string) bool {
	_, err := time.Parse(db.DateFormat, date)
	return err == nil
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
