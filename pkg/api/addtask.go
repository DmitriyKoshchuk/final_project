package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DmitriyKoshchuk/final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var t db.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if t.Date == "" || t.Title == "" {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "Date and title are required"})
		return
	}

	id, err := db.AddTask(&t)
	if err != nil {
		writeJsonWithCode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if task.Date == "" {
		task.Date = today.Format(DateFormat)
	}

	t, err := time.ParseInLocation(DateFormat, task.Date, now.Location())
	if err != nil {
		return err
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if !t.After(today) {
		if task.Repeat == "" {
			task.Date = today.Format(DateFormat)
		} else {
			if t.Before(today) {
				task.Date = next
			}
		}
	}

	return nil
}

func writeJson(w http.ResponseWriter, data interface{}) {
	writeJsonWithCode(w, http.StatusOK, data)
}

func writeJsonWithCode(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if data == nil {
		w.Write([]byte("{}"))
		return
	}

	_ = json.NewEncoder(w).Encode(data)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonWithCode(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJsonWithCode(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	if task.ID == "" {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	if task.Title == "" {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJsonWithCode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJsonWithCode(w, http.StatusOK, nil)
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJsonWithCode(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat != "" {
		nextDate, err := NextDate(taskDateToTime(task.Date), task.Date, task.Repeat)
		if err != nil {
			writeJsonWithCode(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка расчёта следующей даты"})
			return
		}

		err = db.UpdateDate(nextDate, task.ID)
		if err != nil {
			writeJsonWithCode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	} else {
		err = db.DeleteTask(task.ID)
		if err != nil {
			writeJsonWithCode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJsonWithCode(w, http.StatusOK, nil)
}

func taskDateToTime(dateStr string) time.Time {
	t, _ := time.Parse("20060102", dateStr)
	return t
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "Missing id"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJsonWithCode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJsonWithCode(w, http.StatusOK, nil)
}
