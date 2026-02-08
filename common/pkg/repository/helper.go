package repository

import (
	"fmt"
	"time"
)

// ParseSQLiteTime parses datetime strings returned by SQLite (TEXT columns) into time.Time.
// Use this when scanning TEXT datetime columns (e.g. from datetime('now')) into time.Time.
func ParseSQLiteTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid datetime %q", s)
}
