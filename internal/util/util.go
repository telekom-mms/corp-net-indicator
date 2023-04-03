package util

import "time"

const DATE_TIME_FORMAT = "02.01.2006 15:04:05"
const DefaultValue = "-"

// formats date values accordingly "dd.dd.YYYY hh:mm:ss" date format
func FormatDate(t int64) string {
	if t <= 0 {
		return DefaultValue
	}
	return time.Unix(t, 0).Local().Format(DATE_TIME_FORMAT)
}

// formats values -> empty strings are display as "-"
func FormatValue(v string) string {
	if v == "" {
		return DefaultValue
	}
	return v
}
