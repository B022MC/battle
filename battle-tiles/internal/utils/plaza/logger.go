package plaza

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
)

// logger 是本包通用日志器；默认输出到标准输出。
var logger = log.NewHelper(log.With(log.NewStdLogger(os.Stdout), "module", "utils/plaza"))

// SetLogger 允许外部注入统一的 logger。
func SetLogger(l log.Logger) {
	if l == nil {
		return
	}
	logger = log.NewHelper(log.With(l, "module", "utils/plaza"))
}
