package models

import "time"

type BuoyObservation struct {
	BuoyId     string
	ObservedAt time.Time
}
