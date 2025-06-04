package api

import (
	"go1f/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		writeError(w, err)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJson(w, TasksResp{Tasks: tasks})
}
