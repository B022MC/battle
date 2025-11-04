package timeutil

import "time"

type LocalTime time.Time

var timeFormats = []string{time.RFC3339, time.DateTime, "2006-01-02 15", "2006-01-02 15:04", time.DateOnly, time.TimeOnly}

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	// 空值不进行解析
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return nil
	}

	str := string(data)
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	var now time.Time

	for _, format := range timeFormats {

		if now, err = time.ParseInLocation(format, str, time.Local); err == nil {
			*t = LocalTime(now)
			return
		}

	}

	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {

	b := make([]byte, 0, len(time.DateTime)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, time.DateTime)
	b = append(b, '"')
	return b, nil
}
