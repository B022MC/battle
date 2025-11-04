package timeutil

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cast"
)

// CSTLayout China Standard Time Layout
const (
	DiagonalDateLayout  = "2006/01/02"
	CSTDateLayout       = "2006-01-02"
	CSTMonthLayout      = "2006-01"
	CSTMonthDayLayout   = "01-02"
	CSTLayout           = "2006-01-02 15:04:05"
	CSTHourMinuteLayout = "15:04"

	CSTDateYearLayout   = "2006"
	CSTDateMonthLayout  = "200601"
	CSTDateDayLayout    = "20060102"
	CSTDateHourLayout   = "2006010215"
	CSTDateMinuteLayout = "20060102150405"

	StepTypeYear   = "year"
	StepTypeMonth  = "month"
	StepTypeDay    = "day"
	StepTypeHour   = "hour"
	StepTypeMinute = "minute"
	StepTypeSecond = "second"
)

// CSTLayoutString 格式化时间
// 返回 "2006-01-02 15:04:05" 格式的时间
func CSTLayoutString() string {
	ts := time.Now()
	return ts.Format(CSTLayout)
}

// CSTLayoutString 格式化时间
// 返回 "2006-01-02" 格式的时间
func CSTDateLayoutString() string {
	ts := time.Now()
	return ts.Format(CSTDateLayout)
}

// UnixStringToCSTDate 将 Unix 时间戳转换为 "2006-01-02" 格式的日期字符串
func UnixStringToCSTDate(timestampStr string) (string, error) {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", err
	}
	// 将 Unix 时间戳转换为 CST 格式化日期
	return time.UnixMilli(timestamp).Format(CSTDateLayout), nil
}

// UnixStringToCST 将 Unix 时间戳转换为 "2006-01-02 15:04:05" 格式的日期字符串
func UnixStringToCST(timestampStr string) (string, error) {
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return "", err
	}
	// 将 Unix 时间戳转换为 CST 格式化日期
	return time.UnixMilli(timestamp).Format(CSTLayout), nil
}

// ConvertToCustomFormat 将 2024-09-02T06:00:00+08:00 格式的时间字符串转换为 "2006-01-02 15:04:05" 格式的日期字符串
func ConvertToCustomFormat(datetimeStr string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, datetimeStr)
	if err != nil {
		return "", err
	}
	return parsedTime.Format("2006-01-02 15:04:05"), nil
}

// ParseLocalTime 转换成 *LocalTime
func ParseLocalTime(timeStr string) (*LocalTime, error) {
	formats := []string{
		CSTLayout,    // 格式：2024-11-21 10:00:00
		time.RFC3339, // 格式：2024-11-21T10:00:00+08:00
	}

	var parsedTime time.Time
	var err error

	for _, format := range formats {
		parsedTime, err = time.ParseInLocation(format, timeStr, time.Local)
		if err == nil {
			localTime := LocalTime(parsedTime)
			return &localTime, nil
		}
	}

	return nil, fmt.Errorf("时间格式化错误: %w", err)
}

// ParseRFC3339DateLayout 格式化时间
func ParseRFC3339DateLayout(date string) (time.Time, error) {
	return time.Parse(time.RFC3339, date)
}

// ParseCSTDateLayout 格式化时间
func ParseCSTDateLayout(date string) (time.Time, error) {
	return time.ParseInLocation(CSTDateLayout, date, time.Local)
}

// ParseCSTLayout 格式化时间
func ParseCSTLayout(datetime string) (time.Time, error) {
	return time.ParseInLocation(CSTLayout, datetime, time.Local)
}

func ParseCSTAllLayout(datetime string) (time.Time, error) {
	if datetime == "" {
		return time.Time{}, errors.New("时间不能为空")
	}
	if len(datetime) <= 10 {
		return ParseCSTDateLayout(datetime)
	}
	return ParseCSTLayout(datetime)
}

// ParseCSTLayoutByInt64 格式化时间
func ParseCSTLayoutByInt64(date int64, layout string) (time.Time, error) {
	return time.ParseInLocation(layout, strconv.FormatInt(date, 10), time.Local)
}

func ToMinuteTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
}

func ReplaceTime5Minute(start, end time.Time) (time.Time, time.Time) {
	startRemainder := start.Minute() % 5
	if startRemainder == 0 {
		start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), 0, 0, time.Local)
	}
	if startRemainder > 0 {
		add := start.Add(time.Duration(0-startRemainder) * time.Minute)
		start = time.Date(add.Year(), add.Month(), add.Day(), add.Hour(), add.Minute(), 0, 0, time.Local)
	}
	endRemainder := end.Minute() % 5
	if endRemainder > 0 {
		add := end.Add(time.Duration(5-endRemainder) * time.Minute)
		end = time.Date(add.Year(), add.Month(), add.Day(), add.Hour(), add.Minute(), 0, 0, time.Local)
	}
	return start, end
}

func generateBetweenVal(start, end int) []int {
	var result []int
	if start > end {
		return result
	}
	if start == end {
		result = append(result, start)
		return result
	}
	result = append(result, start)
	result = append(result, end)
	return result
}

// GenerateIntervalVal 生成指定时间范围内的时间列表，接受 start 和 end 参数，间隔可以设置为小时、分钟或秒
func GenerateIntervalVal(start, end time.Time, step int, stepType string) ([]time.Time, error) {
	var dur time.Duration
	switch stepType {
	case StepTypeHour:
		dur = time.Duration(step) * time.Hour
	case StepTypeMinute:
		dur = time.Duration(step) * time.Minute
	case StepTypeSecond:
		dur = time.Duration(step) * time.Second
	default:
		return nil, errors.New("unsupported stepType")
	}

	// 将 start 和 end 的分钟和秒设置为 0，调整为整点时间
	if start.Minute() != 0 || start.Second() != 0 {
		start = start.Add(time.Hour - time.Duration(start.Minute())*time.Minute - time.Duration(start.Second())*time.Second)
	}
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), 0, 0, 0, end.Location())

	if start.After(end) {
		return nil, errors.New("start time must be before end time")
	}

	var result []time.Time
	for t := start; t.Before(end) || t.Equal(end); t = t.Add(dur) {
		result = append(result, t)
	}

	return result, nil
}

// GenerateIntervalVal4 生成指定时间范围内的时间列表，接受 start 和 end 参数，间隔可以设置为小时、分钟或秒
func GenerateIntervalVal4(start, end time.Time, step int, stepType string) ([]time.Time, error) {
	if step < 1 {
		step = 1
	}

	var duration time.Duration
	switch stepType {
	case StepTypeHour:
		duration = time.Duration(step) * time.Hour
	case StepTypeMinute:
		duration = time.Duration(step) * time.Minute
	case StepTypeSecond:
		duration = time.Duration(step) * time.Second
	default:
		return nil, errors.New("unsupported stepType")
	}

	result := []time.Time{start}

	var firstAligned time.Time
	if start.Minute() == 0 && start.Second() == 0 {
		firstAligned = start
	} else {
		if stepType == StepTypeHour {
			firstAligned = start.Truncate(time.Hour).Add(time.Hour)
		} else {
			firstAligned = start.Add(duration)
		}
	}

	if !firstAligned.After(start) {
		firstAligned = firstAligned.Add(duration)
	}

	for t := firstAligned; t.Before(end); t = t.Add(duration) {
		result = append(result, t)
	}

	return result, nil
}

//func GenerateIntervalVal(start, end time.Time, step int, stepType string) ([]time.Time, error) {
//	var dur time.Duration
//	switch stepType {
//	case StepTypeHour:
//		dur = time.Duration(step) * time.Hour
//	case StepTypeMinute:
//		dur = time.Duration(step) * time.Minute
//	case StepTypeSecond:
//		dur = time.Duration(step) * time.Second
//	default:
//		return nil, errors.New("unsupported stepType")
//	}
//
//	// 将 start 和 end 的分钟和秒设置为 0
//	start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), 0, 0, 0, start.Location())
//	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), 0, 0, 0, end.Location())
//
//	if start.After(end) {
//		return nil, errors.New("start time must be before end time")
//	}
//
//	var result []time.Time
//	for t := start; t.Before(end) || t.Equal(end); t = t.Add(dur) {
//		result = append(result, t)
//	}
//
//	return result, nil
//}
//

func GenerateIntervalVal2(start, end time.Time, step int, stepType string) ([]int64, error) {
	years := 0
	months := 0
	days := 0
	layout := CSTDateDayLayout
	switch stepType {
	case StepTypeYear:
		years = step
		layout = CSTDateYearLayout
	case StepTypeMonth:
		months = step
		layout = CSTDateMonthLayout
	case StepTypeDay:
		days = step
		layout = CSTDateDayLayout
	}

	var result []int64
	if start.After(end) {
		return result, errors.New("start time must be before end time")
	}
	if start.Equal(end) {
		curTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
		date, _ := strconv.ParseInt(curTime.Format(layout), 10, 64)
		result = append(result, date)
		return result, errors.New("start time must be before end time")
	}
	for t := start; t.Before(end); t = t.AddDate(years, months, days) {
		curTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		date, _ := strconv.ParseInt(curTime.Format(layout), 10, 64)
		result = append(result, date)
	}
	return result, nil
}

// 手麻的
func GenerateIntervalVal3(start, end time.Time, step int) ([]int64, error) {
	var result []int64
	if start.After(end) {
		return result, errors.New("开始时间必须在结束时间之前")
	}
	if start.Equal(end) {
		curTime := time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), 0, 0, time.Local)
		minute, _ := strconv.ParseInt(curTime.Format(CSTDateMinuteLayout), 10, 64)
		result = append(result, minute)
		return result, nil
	}
	for t := start; !t.After(end); t = t.Add(time.Duration(step) * time.Minute) {
		curTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
		minute, _ := strconv.ParseInt(curTime.Format(CSTDateMinuteLayout), 10, 64)

		result = append(result, minute)
	}
	return result, nil
}

func ReplaceMonthDayHOurMinuteIntTime(date time.Time) (month, day, hour int, minute string, err error) {
	month, err = strconv.Atoi(date.Format(CSTDateMonthLayout))
	if err != nil {
		return 0, 0, 0, "", err
	}
	day, err = strconv.Atoi(date.Format(CSTDateDayLayout))
	if err != nil {
		return 0, 0, 0, "", err
	}
	hour, err = strconv.Atoi(date.Format(CSTDateHourLayout))
	if err != nil {
		return 0, 0, 0, "", err
	}

	minute = time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), 0, 0, time.Local).Format(CSTDateMinuteLayout)

	return month, day, hour, minute, nil
}

func ReplaceIntTime(start, end time.Time) (months []int, days []int, hours []int, err error) {

	startMonth, err := strconv.Atoi(start.Format(CSTDateMonthLayout))
	endMonth, err := strconv.Atoi(end.Format(CSTDateMonthLayout))
	if err != nil {
		return nil, nil, nil, err
	}
	startDay, err := strconv.Atoi(start.Format(CSTDateDayLayout))
	endDay, err := strconv.Atoi(end.Format(CSTDateDayLayout))
	if err != nil {
		return nil, nil, nil, err
	}
	startHour, err := strconv.Atoi(start.Format(CSTDateHourLayout))
	endHour, err := strconv.Atoi(end.Format(CSTDateHourLayout))
	if err != nil {
		return nil, nil, nil, err
	}
	months = generateBetweenVal(startMonth, endMonth)
	days = generateBetweenVal(startDay, endDay)
	hours = generateBetweenVal(startHour, endHour)
	return months, days, hours, nil
}

// GenerateTimeIntervalsWithoutFullHours 生成从开始时间到结束时间的指定间隔的分钟数列表，不包含整点
// start 和 end 是 起始和结束时间的字符串格式，例如 "2024-10-14 08:00:00"。
// interval 是间隔的分钟数，表示每次增加的时间间隔。
func GenerateTimeIntervalsWithoutFullHours(start, end string, interval int32) ([]string, error) {
	var intervals []string
	startTime, err := time.ParseInLocation(CSTLayout, start, time.Local)
	if err != nil {
		return nil, fmt.Errorf("解析 start 时间错误: %v", err)
	}
	endTime, err := time.ParseInLocation(CSTLayout, end, time.Local)
	if err != nil {
		return nil, fmt.Errorf("解析 end 时间错误: %v", err)
	}

	for t := startTime; t.Before(endTime.Add(time.Duration(1) * time.Minute)); t = t.Add(time.Duration(interval) * time.Minute) {
		// 跳过整点
		if t.Minute() == 0 {
			continue
		}
		intervals = append(intervals, t.Format(CSTLayout))
	}
	return intervals, nil
}

// GenerateTimeIntervalsWithoutFullHoursTime 生成从开始时间到结束时间的指定间隔的时间点列表（不包含整点分钟）
// 参数为 time.Time 类型，interval 是间隔分钟数。
func GenerateTimeIntervalsWithoutFullHoursTime(start, end time.Time, interval int32) []time.Time {
	var intervals []time.Time

	for t := start; t.Before(end.Add(time.Minute)); t = t.Add(time.Duration(interval) * time.Minute) {
		if t.Minute() == 0 {
			continue
		}
		intervals = append(intervals, t)
	}

	return intervals
}

// GenerateTimeListInRange 生成指定时间范围内的时间列表，接受字符串形式的 startTime 和 endTime，间隔以小时为单位
func GenerateTimeListInRange(startTime, endTime time.Time, interval int32) ([]string, error) {
	var timeList []string
	if interval < 1 {
		interval = 1
	}
	duration := time.Duration(interval) * time.Minute
	// 如果开始时间不是整点，调整到下一个整点
	if startTime.Minute() != 0 || startTime.Second() != 0 {
		// 将开始时间的分钟和秒数调整到 0，得到整点时间
		startTime = startTime.Add(time.Hour - time.Duration(startTime.Minute())*time.Minute - time.Duration(startTime.Second())*time.Second)
	}

	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(duration) {
		timeList = append(timeList, t.Format(CSTLayout))
	}

	return timeList, nil
}

// GenerateTimeListInRange2 生成指定时间范围内的时间列表， startTime 和 endTime，间隔以小时为单位
func GenerateTimeListInRange2(startTime, endTime time.Time, interval int32) ([]string, error) {
	var timeList []string

	if interval < 1 {
		interval = 60 // 默认 60分钟
	}
	duration := time.Duration(interval) * time.Minute

	timeList = append(timeList, startTime.Format(CSTLayout))
	var firstAligned time.Time
	if startTime.Minute() == 0 && startTime.Second() == 0 {
		firstAligned = startTime
	} else {
		firstAligned = startTime.Truncate(time.Hour).Add(time.Hour)
	}
	if !firstAligned.After(startTime) {
		firstAligned = firstAligned.Add(duration)
	}
	for t := firstAligned; t.Before(endTime); t = t.Add(duration) {
		timeList = append(timeList, t.Format(CSTLayout))
	}
	return timeList, nil
}

// 生成从开始时间到结束时间的每五分钟的整点数据
func generateTimeIntervals(start, end time.Time) []time.Time {
	var intervals []time.Time

	for t := start; t.Before(end); t = t.Add(5 * time.Minute) {
		intervals = append(intervals, t)
	}

	return intervals
}

func ToToDayStartAndEndTime(t time.Time) (time.Time, time.Time) {
	startTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	endTime := startTime.AddDate(0, 0, 1).Add(-1 * time.Second)
	return startTime, endTime
}

func ToToMonthStartAndEndTime(t time.Time) (time.Time, time.Time) {
	startTime := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	endTime := startTime.AddDate(0, 1, 0).Add(-1 * time.Second)
	return startTime, endTime
}

func ToStartAndEndTimeByDomain(start, end time.Time) (time.Time, time.Time) {

	startTime := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.Local)
	endTime := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1).Add(-1 * time.Second)

	return startTime, endTime
}

// ToDateDay 将时间转换为 yyyymmdd 格式的整数表示
func ToDateDay(d time.Time) *int32 {
	day := cast.ToInt32(fmt.Sprintf("%d%02d%02d", d.Year(), d.Month(), d.Day()))
	return &day
}

// ToDateMonth 将时间转换为 yyyymm 格式的整数表示
func ToDateMonth(d time.Time) *int32 {
	month := cast.ToInt32(fmt.Sprintf("%d%02d", d.Year(), d.Month()))
	return &month
}

// ToDateYear 将时间转换为 yyyy 格式的整数表示
func ToDateYear(d time.Time) *int32 {
	year := cast.ToInt32(d.Year())
	return &year
}
func ToCSTHourMinute(timeStr string) string {
	parsedTime, err := time.ParseInLocation(CSTLayout, timeStr, time.Local)
	if err != nil {
		return ""
	}
	return parsedTime.Format("15")
}
