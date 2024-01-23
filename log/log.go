package log

import (
	"github.com/charmbracelet/log"
	"os"
	"time"
)

var Default *log.Logger

func init() {
	Default = log.NewWithOptions(os.Stderr, log.Options{
		Level:           log.DebugLevel,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		CallerOffset:    1,
	})
}

// Print prints a log message without levels and colors.
func Print(msg interface{}, keyvals ...interface{}) {
	Default.Print(msg, keyvals...)
}

// Printf prints a log message without levels and colors.
func Printf(format string, args ...interface{}) {
	Default.Printf(format, args...)
}

// Fatal `os.Exit(1)` exit no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatal(msg interface{}, keyvals ...interface{}) {
	Default.Fatal(msg, keyvals...)
}

// Fatalf will `os.Exit(1)` no matter the level of the logger.
// If the logger's level is fatal, error, warn, info or debug
// then it will print the log message too.
func Fatalf(format string, args ...interface{}) {
	Default.Fatalf(format, args...)
}

// Error will print only when logger's Level is error, warn, info or debug.
func Error(msg interface{}, keyvals ...interface{}) {
	Default.Error(msg, keyvals...)
}

// Errorf will print only when logger's Level is error, warn, info or debug.
func Errorf(format string, args ...interface{}) {
	Default.Errorf(format, args...)
}

// Warn will print when logger's Level is warn, info or debug.
func Warn(msg interface{}, keyvals ...interface{}) {
	Default.Warn(msg, keyvals...)
}

// Warnf will print when logger's Level is warn, info or debug.
func Warnf(format string, args ...interface{}) {
	Default.Warnf(format, args...)
}

// Info will print when logger's Level is info or debug.
func Info(msg interface{}, keyvals ...interface{}) {
	Default.Info(msg, keyvals...)
}

// Infof will print when logger's Level is info or debug.
func Infof(format string, args ...interface{}) {
	Default.Infof(format, args...)
}

// Debug will print when logger's Level is debug.
func Debug(msg interface{}, keyvals ...interface{}) {
	Default.Debug(msg, keyvals...)
}

// Debugf will print when logger's Level is debug.
func Debugf(format string, args ...interface{}) {
	Default.Debugf(format, args...)
}
