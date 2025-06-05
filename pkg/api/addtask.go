package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"go1f/pkg/db"
	"net/http"
	"time"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "не удалось разобрать JSON"})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJson(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	nowStr := now.Format(db.DateFormat)
	nowDate, _ := time.Parse(db.DateFormat, nowStr)

	if task.Date == "" {
		task.Date = nowStr
	}

	t, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		return errors.New("дата указана в неверном формате")
	}

	if task.Repeat != "" {
		if t.Before(nowDate) {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return errors.New("неверное правило повторения")
			}
			task.Date = next
		}
	} else {
		if t.Before(nowDate) {
			task.Date = nowStr
		}
	}

	return nil
}
