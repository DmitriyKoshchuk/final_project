package api

import (
	"fmt"
	"go1f/pkg/db"
	"net/http"
	"time"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

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

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("Ошибка вычисления следующей даты"))
			return
		}

		err = db.UpdateDate(id, next)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
