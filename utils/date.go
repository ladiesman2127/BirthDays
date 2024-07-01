package utils

import (
	"strconv"
	"strings"
)

func DateValid(date *string) bool {
	d, m, y, ok := ParseDate(date)
	if !ok {
		return false
	}
	// Check dd/mm/yyyy
	if *d < 1 || *d > 31 || *m < 1 || *m > 12 || *y < 1 {
		return false
	}
	return true
}

func ParseDate(date *string) (*int, *int, *int, bool) {
	splittedData := strings.Split(*date, "-")
	if len(splittedData) < 3 {
		splittedData = strings.Split(*date, "/")
	}
	if len(splittedData) < 3 {
		splittedData = strings.Split(*date, ".")
	}
	if len(splittedData) != 3 {
		return nil, nil, nil, false
	}
	d, err_d := strconv.Atoi(splittedData[0])
	m, err_m := strconv.Atoi(splittedData[1])
	y, err_y := strconv.Atoi(splittedData[2])
	if err_d != nil || err_m != nil || err_y != nil {
		return nil, nil, nil, false
	}
	return &d, &m, &y, true
}
