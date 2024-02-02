package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/rohanraj7316/logger"
	"github.com/rohanraj7316/utils/constants"

	"gorm.io/gorm"
	sl "gorm.io/gorm/logger"
)

const (
	preFixWarn = "[SQL-QUERY]"
	preFixErr  = "[SQL-ERR]"
)

type log struct {
	SlowQueryThreshold time.Duration
	LogLevel           sl.LogLevel
}

func (l *log) LogMode(level sl.LogLevel) sl.Interface {
	return l
}

func (l log) Info(_ context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= sl.Info {
		logger.Info(fmt.Sprintf(msg, args...))
	}
}

func (l log) Warn(_ context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= sl.Warn {
		logger.Warn(fmt.Sprintf(msg, args...))
	}
}

func (l log) Error(_ context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= sl.Error {
		logger.Error(fmt.Sprintf(msg, args...))
	}
}

func (l log) Trace(ctx context.Context, begin time.Time,
	fc func() (sql string, rowsAffected int64), err error) {
	fields := []logger.Field{}

	rId := ctx.Value(constants.REQUEST_ID_PROP)
	fields = append(fields, logger.Field{Key: constants.REQUEST_ID_PROP, Value: rId})

	latency := time.Since(begin).Round(time.Millisecond)
	fields = append(fields, logger.Field{Key: "latency", Value: latency})

	switch {
	case err != nil && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		erStr := fmt.Sprintf("%s | %s | %d", preFixErr, latency.String(), rows)

		fields = append(fields, []logger.Field{
			{Key: "query", Value: sql},
			{Key: "rows", Value: rows},
		}...)

		logger.Error(erStr, fields...)
	case latency >= l.SlowQueryThreshold:
		sql, rows := fc()
		slowStr := fmt.Sprintf("%s | %s | %d", preFixWarn, latency.String(), rows)

		fields = append(fields, []logger.Field{
			{Key: "query", Value: sql},
			{Key: "rows", Value: rows},
		}...)

		logger.Warn(slowStr, fields...)
	}
}
