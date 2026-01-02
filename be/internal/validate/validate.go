package validate

import (
	"regexp"
	"time"
)

var (
	monthRe = regexp.MustCompile(`^\d{4}-\d{2}$`)
	dateRe  = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

func MonthKey(v string) bool {
	if !monthRe.MatchString(v) {
		return false
	}
	_, err := time.Parse("2006-01", v)
	return err == nil
}

func DateKey(v string) bool {
	if !dateRe.MatchString(v) {
		return false
	}
	_, err := time.Parse("2006-01-02", v)
	return err == nil
}
