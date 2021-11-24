package times

import "time"

func Zero() time.Time {
	return time.Time{}
}

func ToYMDHIS(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
