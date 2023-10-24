package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func ConvertMMddStringDateToDate(strDate string) (formattedDate time.Time, err error) {
	dateParts := strings.Split(strDate, "/")
	if len(dateParts) != 2 {
		return formattedDate, errors.New("Invalid date format")
	}
	strFormattedDate := fmt.Sprintf("%02s/%02s", dateParts[0], dateParts[1])

	return time.Parse("01/02", strFormattedDate)
}
