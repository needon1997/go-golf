package golf

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)

// SetLogger sets a new logger.
func SetLogger(l *log.Logger) {
	Logger = l
}

// Trace logs a message at trace level.
func Trace(v ...interface{}) {
	if level <= LevelTrace {
		Logger.Printf("[T] %v\n", v)
	}
}

// Debug logs a message at debug level.
func Debug(v ...interface{}) {
	if level <= LevelDebug {
		Logger.Printf("[D] %v\n", v)
	}
}

// Info logs a message at info level.
func Info(v ...interface{}) {
	if level <= LevelInfo {
		Logger.Printf("[I] %v\n", v)
	}
}

// Warning logs a message at warning level.
func Warn(v ...interface{}) {
	if level <= LevelWarning {
		Logger.Printf("[W] %v\n", v)
	}
}

// Error logs a message at error level.
func Error(v ...interface{}) {
	if level <= LevelError {
		Logger.Printf("[E] %v\n", v)
	}
}

// Critical logs a message at critical level.
func Critical(v ...interface{}) {
	if level <= LevelCritical {
		Logger.Printf("[C] %v\n", v)
	}
}

const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
)

// logLevel controls the global log level used by the logger.
var level = LevelTrace

// LogLevel returns the global log level and can be used in
// own implementations of the logger interface.
func Level() int {
	return level
}

// SetLogLevel sets the global log level used by the simple
// logger.
func SetLevel(l int) {
	level = l
}
