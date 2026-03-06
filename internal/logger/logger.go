package logger

import (
	"fmt"
	"os"
	"strings"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var currentLevel = LevelInfo

func SetLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug":
		currentLevel = LevelDebug
	case "info":
		currentLevel = LevelInfo
	case "warn":
		currentLevel = LevelWarn
	case "error":
		currentLevel = LevelError
	default:
		return fmt.Errorf("invalid log level: %s (valid: debug, info, warn, error)", level)
	}
	return nil
}

func Debug(format string, args ...interface{}) {
	if currentLevel <= LevelDebug {
		fmt.Fprintf(os.Stderr, "DEBUG: "+format+"\n", args...)
	}
}

func Info(format string, args ...interface{}) {
	if currentLevel <= LevelInfo {
		fmt.Fprintf(os.Stderr, "INFO: "+format+"\n", args...)
	}
}

func Warn(format string, args ...interface{}) {
	if currentLevel <= LevelWarn {
		fmt.Fprintf(os.Stderr, "WARN: "+format+"\n", args...)
	}
}

func Error(format string, args ...interface{}) {
	if currentLevel <= LevelError {
		fmt.Fprintf(os.Stderr, "ERROR: "+format+"\n", args...)
	}
}
