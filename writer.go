package go_db

import (
	"go.uber.org/zap"
	"strings"
)

type ZapWriter struct {
	logger *zap.SugaredLogger
}

// 问题 1 infoStr([info])这些是可以改的 2 外部日志和数据库日志两者按最高的来显示
func (w ZapWriter) Printf(format string, args ...interface{}) {
	if w.logger == nil {
		w.logger.Error("错误， 未设置logger")
		return
	}
	if strings.Contains(format, "[info]") {
		w.logger.Infof(format, args)
		return
	} else if strings.Contains(format, "[warn]") {
		w.logger.Warnf(format, args)
		return
	} else if strings.Contains(format, "[error]") {
		w.logger.Errorf(format, args)
		return
	} else {
		w.logger.Infof(format, args)
	}
}
