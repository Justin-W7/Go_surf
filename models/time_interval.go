package models

import (
	"fmt"
	"strings"
	"time"
)

// TimeInterval represents a time range in the format "start/duration".
type TimeInterval struct {
	Start    time.Time
	Duration string
}

// UnmarshalJSON parses a JSON string in the format "start/duration" into
// a TimeInterval.
//
// The JSON value should be a string, e.g. "2025-12-27T15:00:00Z/PT1H".
//
// Parameters:
//   - b: The raw JSON byte slice to parse.
//
// Returns:
//   - error: An error if the string is not in the expected format or
//     if parsing the start time fails.
func (ti *TimeInterval) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid validTime format: %s", s)
	}

	start, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return err
	}

	ti.Start = start
	ti.Duration = parts[1]
	return nil
}
