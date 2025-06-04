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
		writeJSON(w, map[string]string{"error": "не удалось разобрать JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка при добавлении задачи"})
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func checkDate(task *db.Task) error {
	const dateFormat = "20060102"
	now := time.Now()
	nowStr := now.Format(dateFormat)
	nowDate, _ := time.Parse(dateFormat, nowStr)

	if task.Date == "" {
		task.Date = nowStr
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("дата указана в неверном формате")
	}

	if task.Repeat != "" {
		if t.Before(nowDate) {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return errors.New("неверное правило повторения")
			}
			fmt.Printf("Repeat task date before now, updated %s -> %s\n", task.Date, next)
			task.Date = next
		}
	} else {
		if t.Before(nowDate) {
			fmt.Printf("One-time task date before now, updated %s -> %s\n", task.Date, nowStr)
			task.Date = nowStr
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(data)
}
