package db

import (
	"fmt"
	"strconv"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query(
		`SELECT id, date, title, comment, repeat 
         FROM scheduler 
         ORDER BY date 
         LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var id int64
		t := new(Task)
		err := rows.Scan(&id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		t.ID = fmt.Sprintf("%d", id)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var t Task
	var idNum int64
	err := DB.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id).
		Scan(&idNum, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}
	t.ID = fmt.Sprintf("%d", idNum)
	return &t, nil
}

func UpdateTask(task *Task) error {
	idNum, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("Неверный ID")
	}

	res, err := DB.Exec(
		`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, idNum)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

func DeleteTask(id string) error {
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("Неверный ID")
	}

	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, idNum)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

func UpdateDate(id string, date string) error {
	idNum, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("Неверный ID")
	}

	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, date, idNum)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}
