package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once           sync.Once
	loggerInstance *zap.Logger
)

type customCore struct {
	file *os.File
}

func (c *customCore) Enabled(l zapcore.Level) bool             { return true }
func (c *customCore) With(fields []zapcore.Field) zapcore.Core { return c }
func (c *customCore) Check(e zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(e.Level) {
		return ce.AddCore(e, c)
	}
	return ce
}

func (c *customCore) Write(e zapcore.Entry, fields []zapcore.Field) error {
	timestamp := e.Time.Format("2006-01-02 15:04:05")
	caller := e.Caller.TrimmedPath()
	if caller == "" {
		caller = "unknown"
	}
	loggerName := e.LoggerName
	if loggerName == "" {
		loggerName = "auto"
	}

	levelStr := e.Level.CapitalString()

	msg := fmt.Sprintf("[%s][%s][%s][%s] %s\n",
		timestamp,
		levelStr,
		loggerName,
		caller,
		e.Message,
	)

	_, err := c.file.WriteString(msg)
	return err
}

func (c *customCore) Sync() error {
	return c.file.Sync()
}

func Get() *zap.Logger {
	once.Do(func() {
		f, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

		core := &customCore{file: f}
		loggerInstance = zap.New(core, zap.AddCaller())
	})

	return loggerInstance
}

func Fatal(message string) {
	Get().WithOptions(zap.AddCallerSkip(1)).Fatal(message)
}

func Error(message string) {
	Get().WithOptions(zap.AddCallerSkip(1)).Error(message)
}

func Warning(message string) {
	Get().WithOptions(zap.AddCallerSkip(1)).Warn(message)
}

func Debug(message string) {
	Get().WithOptions(zap.AddCallerSkip(1)).Debug(message)
}

func Info(message string) {
	Get().WithOptions(zap.AddCallerSkip(1)).Info(message)
}
