package timezone

import "time"

func ToISO8601(value time.Time, timeZone time.Duration) string {
	return value.UTC().Add(timeZone * time.Hour).Format(time.RFC3339Nano)
}

func NowISO8601(timeZone time.Duration) string {
	return ToISO8601(time.Now(), timeZone)
}
