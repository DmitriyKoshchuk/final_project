package api

import (
	"go1f/pkg/db"
	"net/http"
)

const TasksLimit = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	tasks, err := db.Tasks(TasksLimit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJson(w, http.StatusOK, TasksResp{Tasks: tasks})
}
