package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DmitriyKoshchuk/final_project/pkg/db"
)

const dateformat = "20060102"

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Декодируем JSON из тела запроса
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "invalid JSON"})
		return
	}

	// Проверка обязательного поля title
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title is required"})
		return
	}

	// Проверяем и корректируем дату задачи
	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "failed to add task"})
		return
	}

	// Возвращаем id созданной задачи
	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func checkDate(task *db.Task) error {
	// Сегодняшняя дата без времени
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Если дата пустая — берем сегодня
	if task.Date == "" {
		task.Date = today.Format(dateformat)
	}

	// Парсим дату задачи
	t, err := time.ParseInLocation(dateformat, task.Date, now.Location())
	if err != nil {
		return err
	}

	// Если есть правило повторения, проверяем его и получаем следующую дату
	var next string
	if task.Repeat != "" {
		next, err = NextDate(today, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	// Если дата задачи в прошлом относительно сегодня
	if !t.After(today) {
		if task.Repeat == "" {
			task.Date = today.Format(dateformat) // берем сегодня
		} else {
			// Если дата раньше сегодняшней — берем следующую дату повторения
			if t.Before(today) {
				task.Date = next
			}
			// Если дата равна сегодняшней — оставляем без изменений
		}
	}

	return nil
}

// Вспомогательная функция для отправки JSON-ответа
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "invalid JSON"})
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "id is required"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title is required"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, struct{}{}) // Пустой JSON
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	// Если задача повторяющаяся
	if task.Repeat != "" {
		// Считаем следующую дату
		nextDate, err := NextDate(taskDateToTime(task.Date), task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "Ошибка расчёта следующей даты"})
			return
		}

		// Обновляем дату задачи
		err = db.UpdateDate(nextDate, task.ID)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}

	} else {
		// Если одноразовая — удаляем
		err = db.DeleteTask(task.ID)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJson(w, struct{}{}) // пустой JSON в ответ
}

// Вспомогательная функция для преобразования строки даты "20060102" в time.Time
func taskDateToTime(dateStr string) time.Time {
	t, _ := time.Parse("20060102", dateStr)
	return t
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, struct{}{})
}
