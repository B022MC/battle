package timex

import "time"

func ReplaceTimeNear5Minute(t time.Time) time.Time {
	startRemainder := t.Minute() % 5
	if startRemainder == 0 {
		t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
	}
	if startRemainder > 0 {
		add := t.Add(time.Duration(0-startRemainder) * time.Minute)
		t = time.Date(add.Year(), add.Month(), add.Day(), add.Hour(), add.Minute(), 0, 0, time.Local)
	}
	return t
}
