package arpscan

import "time"

type mac struct {
	Address  string
	FoundAt  time.Time
	LastTime time.Time
}
