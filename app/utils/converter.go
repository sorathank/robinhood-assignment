package utils

import (
	"errors"
	"strconv"
)

func StringToPositiveInt(s string) (int64, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	if i <= int64(0) {
		return 0, errors.New("number is not positive")
	}
	return i, nil
}
