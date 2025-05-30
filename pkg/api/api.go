package api

import (
	"encoding/json"
	"net/http"
	"os"
)

var password string

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path == "/api/task/done" {
			doneTaskHandler(w, r)
		} else {
			addTaskHandler(w, r)
		}
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPasswordFromEnv() string {
	return os.Getenv("TODO_PASSWORD")
}

func Init() {
	password = getPasswordFromEnv()

	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", authMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", authMiddleware(doneTaskHandler))
	http.HandleFunc("/nextdate", nextDayHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJsonWithCode(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if password == "" || req.Password != password {
		writeJsonWithCode(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	token, err := generateToken()
	if err != nil {
		writeJsonWithCode(w, http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		return
	}

	writeJsonWithCode(w, http.StatusOK, map[string]string{"token": token})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if password == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || !validateToken(cookie.Value) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
