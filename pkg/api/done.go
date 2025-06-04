package api

import (
	"go1f/pkg/db"
	"net/http"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "Ошибка вычисления следующей даты"})
			return
		}

		err = db.UpdateDate(id, next)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJson(w, map[string]string{})
}
