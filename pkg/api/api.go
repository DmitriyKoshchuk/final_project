package api

import (
	"encoding/json"
	"net/http"
	"os"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Проверяем путь
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

func Init() {
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", authMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", authMiddleware(doneTaskHandler))
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, map[string]string{"error": "Invalid request"})
		return
	}

	envPass := os.Getenv("TODO_PASSWORD")
	if envPass == "" || req.Password != envPass {
		writeJson(w, map[string]string{"error": "Неверный пароль"})
		return
	}

	token, err := generateToken()
	if err != nil {
		writeJson(w, map[string]string{"error": "Failed to generate token"})
		return
	}

	writeJson(w, map[string]string{"token": token})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("TODO_PASSWORD") == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil || !validateToken(cookie.Value) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
