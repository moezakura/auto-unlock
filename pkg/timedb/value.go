package timedb

import "time"

type value struct {
	RawValue interface{}
	Deadline time.Time
}
