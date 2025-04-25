package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// ParseQueryInt extracts an integer from query parameters, falling back to a default if missing or invalid.
func ParseQueryInt(c *gin.Context, key string, defaultValue int) int {
	valStr := c.DefaultQuery(key, "")

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}

	return val
}

// ParseQueryTime extracts a time.Time from a query param if present and valid.
func ParseQueryTime(c *gin.Context, key string) *time.Time {
	timeStr := c.Query(key)
	if timeStr == "" {
		return nil
	}

	// try parsing with "YYYY-MM-DD" format
	t, err := time.Parse("2006-01-02", timeStr)
	if err == nil {
		return &t
	}

	// if the first format fails, try parsing with RFC3339
	t, err = time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return &t
	}

	if &t != nil {
		// adjust DateFrom to the beginning of the day if it's not nil
		if key == "from" {
			year, month, day := t.Date()
			startOfDay := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
			t = startOfDay

			return &t
		}

		// adjust DateTo to the end of the day if it's not nil
		if key == "to" {
			year, month, day := t.Date()
			endOfDay := time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
			t = endOfDay

			return &t
		}
	}

	return &t
}

// Max returns the greater of two ints.
func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Min returns the smaller of two ints.
func Min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
