package api

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const formatDate = "20060102" // Format for date in YYYYMMDD format

// nextDayHandler handles the /api/nextdate endpoint.
// It expects 'now', 'date', and 'repeat' parameters in the request.
// 'now' is the current date in YYYYMMDD format, 'date' is the start date in YYYYMMDD format,
// and 'repeat' is the repeat pattern (e.g., "d 7" for every 7 days, "y" for yearly).
// It returns the next date based on the repeat pattern.
// If 'now' is not provided, it defaults to the current date.
// If 'date' or 'repeat' is empty, it returns an error.
// If the date format is invalid, it returns an error.
// The response is in YYYYMMDD format.
// Example request: /api/nextdate?now=20240126&date=20240113&repeat=d 7
// Example response: 20240127
func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	now := r.FormValue("now")
	if now == "" {
		now = time.Now().Format(formatDate)
	}
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	nowTime, err := time.Parse(formatDate, now)
	if err != nil {
		http.Error(w, "Invalid 'now' date format, expected YYYYMMDD", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write([]byte(nextDate))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// NextDate calculates the next date based on the provided start date and repeat pattern.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if dstart == "" {
		return "", fmt.Errorf("data cannot be empty")
	}
	if repeat == "" {
		return "", fmt.Errorf("repeat cannot be empty")
	}
	date, err := time.Parse(formatDate, dstart)
	if err != nil {
		return "", err
	}
	date, err = changeRepeat(date, now, repeat) // Change the date based on the repeat pattern
	if err != nil {
		return "", err
	}
	return date.Format(formatDate), nil
}

// changeRepeat modifies the date based on the repeat pattern.
func changeRepeat(date, now time.Time, repeat string) (time.Time, error) {
	arr := strings.Split(repeat, " ")
	switch arr[0] {
	case "d":
		if len(arr) != 2 {
			return date, fmt.Errorf("invalid repeat format for daily: %s", repeat)
		}
		interval, err := strconv.Atoi(arr[1])
		if err != nil {
			return date, err
		}
		if interval <= 0 || interval >= 400 {
			return date, fmt.Errorf("invalid interval for daily: %d", interval)
		}
		return addDay(date, now, interval), nil
	case "y":
		if len(arr) >= 2 {
			return date, fmt.Errorf("invalid repeat format for yearly: %s", repeat)
		}
		return addYear(date, now, 1), nil
	case "w":
		if len(arr) < 2 {
			return date, fmt.Errorf("the number of parameters cannot be less than two for days of the week: %s", repeat)
		}
		return addWeek(date, now, arr[1])
	case "m":
		if len(arr) < 2 {
			return date, fmt.Errorf("invalid repeat format for monthly: %s", repeat)
		}
		return addMonth(date, now, arr[1:])
	default:
		return date, fmt.Errorf("unsupported repeat type: %s", arr[0])
	}
}

// afterNow checks if the date is after the current time.
func afterNow(date, now time.Time) bool {
	dy, dm, dd := date.Date()
	ny, nm, nd := now.Date()

	if dy != ny {
		return dy > ny
	}

	if dm != nm {
		return dm > nm
	}

	return dd > nd
}

// addDay adds the specified interval of days to the date until it is after now.
// It returns the modified date.
func addDay(date, now time.Time, interval int) time.Time {
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}
	return date
}

// addYear adds the specified interval of years to the date until it is after now.
// It returns the modified date.
func addYear(date, now time.Time, interval int) time.Time {
	for {
		date = date.AddDate(interval, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date
}

func addWeek(date, now time.Time, interval string) (time.Time, error) {
	// If the date is before now, set it to now
	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			break
		}
	}

	daysOfWeek := make(map[int]bool)
	interval = strings.TrimSpace(interval) // Trim whitespace from the interval string
	days := strings.Split(interval, ",")   // Split the interval by comma to get individual days

	// Iterate over the provided days of the week
	for _, val := range days {
		day, err := strconv.Atoi(val)
		if err != nil || day < 1 || day > 7 {
			return date, fmt.Errorf("invalid day value: %s", val)
		}
		day = day % 7          // Normalize the day to 0-6 range (1-7 to 0-6)
		daysOfWeek[day] = true // Mark the day as set
	}
	today := int(date.Weekday())
	for i := 0; i < 7; i++ {
		nextDay := (today + i) % 7 // Calculate the next day in the week
		if daysOfWeek[nextDay] {   // Check if this day is in the provided days of the week
			return date.AddDate(0, 0, i), nil // Add the number of days to reach this day
		}
	}
	return date, fmt.Errorf("no valid days of the week provided in interval: %s", interval)
}

func addMonth(date, now time.Time, interval []string) (time.Time, error) {
	// If the date is before now, set it to now
	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			break
		}
	}

	daysOfMonth := strings.Split(strings.TrimSpace(interval[0]), ",")
	var months []string
	// If the interval has more than one part, the second part contains the months
	if len(interval) > 1 {
		months = strings.Split(strings.TrimSpace(interval[1]), ",")
	}
	if len(interval) == 1 {
		months = []string{} // Default to empty if no months are provided
	}

	monthMap := map[int]bool{}
	if len(months) == 0 {
		for i := 1; i <= 12; i++ {
			monthMap[i] = true
		}
	}
	if len(months) > 0 {
		for _, m := range months {
			m, err := strconv.Atoi(m)
			if err != nil || m < 1 || m > 12 {
				return time.Time{}, fmt.Errorf("invalid month: %d", m)
			}
			monthMap[m] = true
		}
	}

	for i := 0; i < 24; i++ {
		nextMonth := date.AddDate(0, i, 0)
		monthNum := int(nextMonth.Month())
		if !monthMap[monthNum] {
			continue
		}
		lastDay := lastDayOfMonth(nextMonth)

		var candidates []int
		for _, d := range daysOfMonth {
			day, err := strconv.Atoi(d)
			if err != nil || day < -2 || day > 31 || day == 0 {
				return date, fmt.Errorf("invalid day: %s", d)
			}
			switch {
			case day > 0 && day <= lastDay:
				candidates = append(candidates, day)
			case day == -1:
				candidates = append(candidates, lastDay)
			case day == -2 && lastDay > 1:
				candidates = append(candidates, lastDay-1)
			}
		}

		sort.Ints(candidates) // Sort candidates to find the next valid date

		for _, d := range candidates {
			candidateDate := time.Date(nextMonth.Year(), nextMonth.Month(), d, 0, 0, 0, 0, date.Location())
			if candidateDate.After(date) {
				return candidateDate, nil
			}
		}
	}
	return date, fmt.Errorf("cannot find suitable date for given rules")
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}
