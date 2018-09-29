package util

import (
	"strings"
	"time"
)

func GetDateRange(from, to string) (time.Time, time.Time) {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	if from == "" || to == "" {
		if from == "" {
			return getRangeByOneDate(to)
		} else {
			return getRangeByOneDate(from)
		}
	} else {
		return getDateRangeWithTime(from, to)
	}

}

func getRangeByOneDate(date string) (time.Time, time.Time) {
	return getDateWitTime(date, true), getDateWitTime(date, false)
}

func getDateRangeWithTime(from, to string) (time.Time, time.Time) {
	return getDateWitTime(from, true), getDateWitTime(to, false)
}

func getDateWitTime(date string, isFirst bool) time.Time {
	layout := "2006-01-02T15:04:05"
	var normalized string
	if isFirst {
		normalized = date + "T00:00:00"
	} else {
		normalized = date + "T23:59:59"
	}
	out, _ := time.Parse(layout, normalized)
	return out
}
