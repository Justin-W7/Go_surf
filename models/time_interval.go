package models

import (
	"fmt"
	"strings"
	"time"
)

type TimeInterval struct {
	Start    time.Time
	Duration string
}

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
