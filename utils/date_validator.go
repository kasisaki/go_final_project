package utils

import (
	"go_final_project/constants"
	"time"
)

func ValidateDate(dateString string) (bool, string) {
	date, err := time.Parse("02.01.2006", dateString)
	if err != nil {
		date, err = time.Parse("2006.01.02", dateString)
	}
	if err != nil {
		date, err = time.Parse(constants.DateLayout, dateString)
	}
	return err == nil, date.Format(constants.DateLayout)
}
