package main

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	DateLayout = "20060102"
)

func NextDate(now time.Time, date string, repeat string) (newDate string, err error) {
	parameterError := errors.New("wrong repeat parameter")
	parameterSet := strings.Split(repeat, " ")
	taskDate, err := time.Parse(DateLayout, date)
	if err != nil {
		return "", err
	}

	lenP := len(parameterSet)
	if lenP == 0 {
		return "", nil
	}
	dates := make(map[time.Time]bool)
	switch parameterSet[0] {
	// this case will add a required number of days but no more than 400
	case "d":
		if lenP != 2 {
			return "", parameterError
		}
		delay, formatErr := strconv.Atoi(parameterSet[1])
		if formatErr != nil {
			return "", parameterError
		}
		if delay > 400 || delay < 1 {
			return "", parameterError
		}
		taskDate = taskDate.AddDate(0, 0, delay)
		for taskDate.Before(now) || taskDate.Equal(now) {
			taskDate = taskDate.AddDate(0, 0, delay)
		}
		return taskDate.Format(DateLayout), nil
	// this case will add a year to the date
	case "y":
		if lenP != 1 {
			return "", parameterError
		}
		taskDate = taskDate.AddDate(1, 0, 0)
		for now.After(taskDate) || now.Equal(taskDate) {
			taskDate = taskDate.AddDate(1, 0, 0)
		}
		return taskDate.Format(DateLayout), nil
	case "w":
		if taskDate.Before(now) {
			taskDate = now
		}
		if lenP != 2 {
			return "", parameterError
		}
		days := strings.Split(parameterSet[1], ",")
		if len(days) > 7 {
			return "", parameterError
		}

		weekDays := make([]int, len(days))
		for i, day := range days {
			weekDays[i], err = strconv.Atoi(day)
			if err != nil {
				return "", err
			}
			if weekDays[i] > 7 || weekDays[i] < 0 {
				return "", parameterError
			}
			if weekDays[i] == 7 {
				weekDays[i] = 0
			}
		}
		slices.Sort(weekDays)
		for _, day := range weekDays {
			dates[nextWeeklyDate(taskDate, time.Weekday(day))] = true
		}
		return findEarliestDate(dates).Format(DateLayout), nil

	case "m":
		if taskDate.Before(now) {
			taskDate = now
		}
		if lenP == 2 {
			// ежемесячно
			strSet := strings.Split(parameterSet[1], ",")
			days := make([]int, len(strSet))
			for i, str := range strSet {
				days[i], err = strconv.Atoi(str)
				if err != nil {
					return "", err
				}
			}
			dates := make(map[time.Time]bool)
			for _, day := range days {
				if day < 0 {
					if day < -2 {
						return "", parameterError
					}
					y, m, _ := taskDate.Date()
					_, _, dayX := time.Date(y, m+1, day+1, 0, 0, 0, 0, taskDate.Location()).Date()
					day = dayX
				}
				if day > 31 {
					return "", parameterError
				}
				dates[nextMonthlyDate(taskDate, 0, day, false)] = true
			}
			return findEarliestDate(dates).Format(DateLayout), nil
		} else if lenP == 3 {
			// с указанием конкретного месяца
			strMonth := strings.Split(parameterSet[2], ",")
			strDays := strings.Split(parameterSet[1], ",")
			months := make([]int, len(strMonth))
			days := make([]int, len(strDays))
			for i, str := range strMonth {
				months[i], err = strconv.Atoi(str)
				if err != nil {
					return "", err
				}
				if months[i] > 12 {
					return "", parameterError
				}
			}
			for i, str := range strDays {
				days[i], err = strconv.Atoi(str)
				if err != nil {
					return "", err
				}
				if days[i] > 31 {
					return "", parameterError
				}
			}
			for _, month := range months {
				for _, day := range days {
					dates[nextMonthlyDate(taskDate, time.Month(month), day, true)] = true
				}
			}
			return findEarliestDate(dates).Format(DateLayout), nil
		}
		return "", parameterError
	default:
		return "", parameterError
	}
}

func nextWeeklyDate(t time.Time, weekday time.Weekday) time.Time {
	days := int(7 + (weekday-t.Weekday())%7)
	y, m, d := t.AddDate(0, 0, days).Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func nextMonthlyDate(t time.Time, month time.Month, day int, monthSpecific bool) time.Time {
	y, m, d := t.Date()
	if monthSpecific {
		if m > month {
			y++
			m = month //return time.Date(y+1, month, d, 0, 0, 0, 0, t.Location())
		} else if m == month {
			if d < day {
				m = month
				// return time.Date(y, month, day, 0, 0, 0, 0, t.Location())
			} else {
				y++
				// return time.Date(y+1, month, d, 0, 0, 0, 0, t.Location())
			}
		} else {
			m = month
			// return time.Date(y, month, day, 0, 0, 0, 0, t.Location())
		}
	} else {
		if d >= day {
			m++
			// return time.Date(y, m+1, day, 0, 0, 0, 0, t.Location())
		}
	}
	for DaysInMonth(y, m) < day {
		m++
	}
	return time.Date(y, m, day, 0, 0, 0, 0, t.Location())
}

func findEarliestDate(dates map[time.Time]bool) time.Time {
	someDateInFuture, _ := time.Parse(DateLayout, "30000101")
	for k, v := range dates {
		if v {
			if k.Before(someDateInFuture) {
				someDateInFuture = k
			}
		}
	}
	return someDateInFuture
}

func DaysInMonth(year int, month time.Month) int {
	// Get the first day of the next month
	firstDayOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)

	// Subtract one day to get the last day of the current month
	lastDayOfCurrentMonth := firstDayOfNextMonth.AddDate(0, 0, -1)

	// Extract the day from the last day of the month
	daysInMonth := lastDayOfCurrentMonth.Day()

	return daysInMonth
}
