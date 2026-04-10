package logger

import (
	"context"
	"log/slog"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	log           *slog.Logger
	slowThreshold time.Duration
	logLevel      gormlogger.LogLevel
}

func NewGormLogger(log *slog.Logger, env string) gormlogger.Interface {
	level := gormlogger.Warn
	if env == "dev" {
		level = gormlogger.Info
	}

	return &GormLogger{
		log:           log,
		slowThreshold: 200 * time.Millisecond,
		logLevel:      level,
	}
}

func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...any) {
	if l.logLevel < gormlogger.Info {
		return
	}
	l.log.InfoContext(ctx, msg, "args", args)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...any) {
	if l.logLevel < gormlogger.Warn {
		return
	}
	l.log.WarnContext(ctx, msg, "args", args)
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...any) {
	if l.logLevel < gormlogger.Error {
		return
	}
	l.log.ErrorContext(ctx, msg, "args", args)
}

func (l *GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	if l.logLevel == gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.logLevel >= gormlogger.Error:
		l.log.ErrorContext(ctx, "gorm query failed",
			"elapsed", elapsed,
			"rows", rows,
			"sql", sql,
			"error", err,
		)

	case elapsed > l.slowThreshold && l.logLevel >= gormlogger.Warn:
		l.log.WarnContext(ctx, "gorm slow query",
			"elapsed", elapsed,
			"rows", rows,
			"sql", sql,
		)

	case l.logLevel >= gormlogger.Info:
		l.log.InfoContext(ctx, "gorm query",
			"elapsed", elapsed,
			"rows", rows,
			"sql", sql,
		)
	}
}
