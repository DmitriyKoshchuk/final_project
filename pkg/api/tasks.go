package api

import (
	"net/http"

	"github.com/DmitriyKoshchuk/final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJsonWithCode(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJsonWithCode(w, http.StatusOK, TasksResp{Tasks: tasks})
}
