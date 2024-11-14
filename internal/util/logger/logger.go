package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

type CustomFormatter struct {
	logrus.TextFormatter
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Custom timestamp format
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// Get the level color
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = 36 // cyan
	case logrus.InfoLevel:
		levelColor = 32 // green
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 37 // white
	}

	// Format the message with colors and structure
	msg := []byte{}

	// Add timestamp
	msg = append(msg, []byte(
		"\x1b[36m"+timestamp+"\x1b[0m")..., // Cyan timestamp
	)
	msg = append(msg, []byte(" | ")...)

	// Add level
	msg = append(msg, []byte(
		"\x1b["+string(rune(levelColor))+"m"+entry.Level.String()+"\x1b[0m")...,
	)
	msg = append(msg, []byte(" | ")...)

	// Add file and line information if available
	if file, ok := entry.Data["file"]; ok {
		if line, ok := entry.Data["line"]; ok {
			msg = append(msg, []byte(
				fmt.Sprintf("\x1b[33m%s:%d\x1b[0m | ", file, line))..., // Yellow for file:line
			)
		}
	}

	// Add other fields if they exist
	if len(entry.Data) > 0 {
		msg = append(msg, []byte("\x1b[35m")...) // Magenta for fields
		for k, v := range entry.Data {
			if k != "file" && k != "line" { // Skip file and line as they're already added
				msg = append(msg, []byte(fmt.Sprintf("%s=%v ", k, v))...)
			}
		}
		msg = append(msg, []byte("\x1b[0m")...)
		msg = append(msg, []byte("| ")...)
	}

	// Add the message
	msg = append(msg, []byte("\x1b[37m"+entry.Message+"\x1b[0m")...) // White message
	msg = append(msg, []byte("\n")...)

	return msg, nil
}

func init() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetReportCaller(true) // Enable caller reporting
	Log.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			DisableLevelTruncation: true,
		},
	})

	// Set log level from environment variable or default to Info
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		Log.SetLevel(logrus.InfoLevel)
	} else {
		Log.SetLevel(logLevel)
	}
}

// Helper functions for structured logging with caller information
func WithFields(fields logrus.Fields) *logrus.Entry {
	_, file, line, _ := runtime.Caller(1)
	fields["file"] = filepath.Base(file)
	fields["line"] = line
	return Log.WithFields(fields)
}

func WithField(key string, value interface{}) *logrus.Entry {
	_, file, line, _ := runtime.Caller(1)
	return Log.WithFields(logrus.Fields{
		key:    value,
		"file": filepath.Base(file),
		"line": line,
	})
}

// New helper functions for direct logging with caller information
func Info(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	}).Info(args...)
}

func Error(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	}).Error(args...)
}

func Warn(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	}).Warn(args...)
}

func Debug(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	}).Debug(args...)
}

func Fatal(args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	Log.WithFields(logrus.Fields{
		"file": filepath.Base(file),
		"line": line,
	}).Fatal(args...)
}
