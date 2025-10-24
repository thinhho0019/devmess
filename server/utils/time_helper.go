package utils

import (
	"strconv"
	"time"
)

func ConvertMilisecondToTime(msStr string) (*time.Time, error) {
	msInt, err := strconv.ParseInt(msStr, 10, 64)
	if err != nil {
		return nil, err
	}
	t := time.Unix(0, msInt*int64(time.Millisecond))
	return &t, nil
}
