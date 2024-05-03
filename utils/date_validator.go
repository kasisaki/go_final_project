package utils

import (
	"time"
)

func ValidateDate(dateString string) (bool, string) {
	date, err := time.Parse("02.01.2006", dateString)
	if err != nil {
		date, err = time.Parse("2006.01.02", dateString)
	}
	if err != nil {
		date, err = time.Parse(DateLayout, dateString)
	}
	return err == nil, date.Format(DateLayout)
}
