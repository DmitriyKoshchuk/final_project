package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	startDate, err := time.ParseInLocation(DateFormat, dstart, time.UTC)
	if err != nil {
		return "", errors.New("invalid start date format")
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("invalid daily format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 || days > 400 {
			return "", errors.New("invalid number of days")
		}

		next := startDate
		for {
			next = next.AddDate(0, 0, days)
			if next.After(now) {
				break
			}
		}
		return next.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("invalid yearly format")
		}

		next := startDate
		for {
			next = next.AddDate(1, 0, 0)

			if startDate.Month() == time.February && startDate.Day() == 29 {
				if !isLeapYear(next.Year()) {
					next = time.Date(next.Year(), time.March, 1, 0, 0, 0, 0, time.UTC)
				}
			}

			if next.After(now) {
				break
			}
		}
		return next.Format(DateFormat), nil

	default:
		return "", errors.New("unsupported repeat format")
	}
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJsonWithCode(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	repeat := r.URL.Query().Get("repeat")
	dstart := r.URL.Query().Get("date")
	nowStr := r.URL.Query().Get("now")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": "invalid now date"})
			return
		}
	}

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		writeJsonWithCode(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJsonWithCode(w, http.StatusOK, map[string]string{"date": next})
}
