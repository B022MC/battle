package utils

const (
	StatisticsMethodByDay   string = "day"   // 日
	StatisticsMethodByMonth string = "month" // 月
	StatisticsMethodByYear  string = "year"  // 年
)
const (
	StepTypeYear   = "year"
	StepTypeMonth  = "month"
	StepTypeDay    = "day"
	StepTypeHour   = "hour"
	StepTypeMinute = "minute"
	StepTypeSecond = "second"
)

func Paginate(page, pageSize int, isPage bool) (limit, offset int) {
	if !isPage {
		return -1, -1
	}
	if page <= 0 {
		page = 1
	}
	switch {
	case pageSize > 1000:
		pageSize = 1000
	case pageSize <= 0:
		pageSize = 10
	}
	offset = (page - 1) * pageSize
	limit = pageSize
	return limit, offset
}
