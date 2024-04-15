package util

import (
	"fmt"
	"time"
)

func MakeCodeByDate(id int64) string {
	return fmt.Sprintf("%d-%02d-%02d-%08d", time.Now().Year(), time.Now().Month(), time.Now().Day(), id)
}

func MakeCodeByYear(id int64) string {
	return fmt.Sprintf("%02d%02d%09d", time.Now().Year()%100, time.Now().Month(), id%100000000)
}
