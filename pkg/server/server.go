package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/DmitriyKoshchuk/final_project/pkg/api"
	"github.com/DmitriyKoshchuk/final_project/pkg/db"
)

const (
	defaultPort   = "7540"
	defaultDBFile = "scheduler.db"
	webDir        = "./web"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	// Инициализация API должна быть ДО инициализации базы данных
	api.Init()

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	// Остальной код...

	// Инициализация базы данных
	if err := db.Init(dbFile); err != nil {
		return fmt.Errorf("failed to initialize DB: %w", err)
	}

	// Обработчик файлов для фронтенда
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	fmt.Printf("Server started at http://localhost:%s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
