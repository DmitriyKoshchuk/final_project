package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat) // <- используем DB
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	var rows *sql.Rows
	var err error

	if search == "" {
		// Без фильтра — просто берем все задачи, отсортированные по дате
		rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`, limit)
	} else {
		// Проверим, не дата ли это (формат 02.01.2006)
		t, errDate := time.Parse("02.01.2006", search)
		if errDate == nil {
			// Это дата — ищем по точному совпадению
			date := t.Format("20060102")
			rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`, date, limit)
		} else {
			// Иначе — ищем подстроку в title или comment (через LIKE)
			likeSearch := "%" + search + "%"
			rows, err = DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`, likeSearch, likeSearch, limit)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("DB query failed: %w", err)
	}
	defer rows.Close()

	tasks := make([]*Task, 0, limit)
	for rows.Next() {
		var id int64
		var t Task

		err = rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, fmt.Errorf("DB scan failed: %w", err)
		}
		t.ID = strconv.FormatInt(id, 10) // <-- тут преобразуем id в строку
		tasks = append(tasks, &t)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)

	t := &Task{}
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}
