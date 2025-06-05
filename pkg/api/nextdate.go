package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go1f/pkg/db"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(db.DateFormat, nowStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("invalid now date"))
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(next)); err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("failed to write response"))
		return
	}
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	start, err := time.Parse(db.DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}

	parts := strings.Fields(repeat)
	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid daily format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid daily interval")
		}
		for {
			start = start.AddDate(0, 0, days)
			if afterNow(start, now) {
				break
			}
		}
		return start.Format(db.DateFormat), nil
	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid yearly format")
		}
		for {
			start = start.AddDate(1, 0, 0)
			if afterNow(start, now) {
				break
			}
		}
		return start.Format(db.DateFormat), nil
	default:
		return "", errors.New("unsupported repeat format")
	}
}

func afterNow(date, now time.Time) bool {
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	if y1 != y2 {
		return y1 > y2
	}
	if m1 != m2 {
		return m1 > m2
	}
	return d1 > d2
}
