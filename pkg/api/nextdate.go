package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	start, err := time.Parse(dateFormat, dstart)
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
		return start.Format(dateFormat), nil
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
		return start.Format(dateFormat), nil
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

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "invalid now date", http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}
