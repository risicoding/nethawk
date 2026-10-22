package logger

import (
	"os"
	"time"

	fp "path/filepath"

	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
)

// Global logger instances accessible from anywhere in your project
var (
	Log *MultiLogger
)

type MultiLogger struct {
	consoleLogger *log.Logger
	fileLogger    *log.Logger
}

// Init initializes the global logger singleton. Call this ONCE in your main.go.
func Init() func() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic("Failed to get user home dir: " + err.Error())
	}

	filePath := fp.Join(homedir, ".local", "state", "nethawk", "nethawk.log")
	// 1. Open file for logging
	fileDir := fp.Dir(filePath)
	err = os.MkdirAll(fileDir, 0755)
	if err != nil {
		panic("Failed to create log dir: " + err.Error())
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	// 2. Setup Console Logger
	cLog := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		CallerOffset:    1,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
	})

	// 3. Customize Lipgloss Styles
	styles := log.DefaultStyles()
	styles.Timestamp = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
	styles.Prefix = lipgloss.NewStyle().Foreground(lipgloss.Color("#F92672")).Bold(true)
	styles.Message = lipgloss.NewStyle().Foreground(lipgloss.Color("#9cfcff"))
	styles.Key = lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E22E"))
	styles.Value = lipgloss.NewStyle().Foreground(lipgloss.Color("#AE81FF"))

	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().SetString("DEBUG").Foreground(lipgloss.Color("#24ff3a")).Bold(true)
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().SetString("INFO ").Foreground(lipgloss.Color("#00f2ff")).Bold(true)
	styles.Levels[log.WarnLevel] = lipgloss.NewStyle().SetString("WARN ").Foreground(lipgloss.Color("#fff200")).Bold(true)
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().SetString("ERROR").Foreground(lipgloss.Color("#f74040")).Bold(true)
	cLog.SetStyles(styles)

	// 4. Setup File Logger (Structured JSON)
	fLog := log.NewWithOptions(file, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.RFC3339,
		Formatter:       log.JSONFormatter,
	})

	// 5. Save to global variable
	Log = &MultiLogger{
		consoleLogger: cLog,
		fileLogger:    fLog,
	}

	// Return a cleanup function to close the file on shutdown
	return func() {
		file.Close()
	}
}

// Global Package level functions so you don't even have to type logger.Log.Info()
func Debug(msg string, keyvals ...any) {
	Log.consoleLogger.Debug(msg, keyvals...)
	Log.fileLogger.Debug(msg, keyvals...)
}
func Info(msg string, keyvals ...any) {
	Log.consoleLogger.Info(msg, keyvals...)
	Log.fileLogger.Info(msg, keyvals...)
}
func Warn(msg string, keyvals ...any) {
	Log.consoleLogger.Warn(msg, keyvals...)
	Log.fileLogger.Warn(msg, keyvals...)
}
func Error(msg string, keyvals ...any) {
	Log.consoleLogger.Error(msg, keyvals...)
	Log.fileLogger.Error(msg, keyvals...)
}

func Debugf(format string, args ...any) {
	Log.consoleLogger.Debugf(format, args...)
	Log.fileLogger.Debugf(format, args...)
}

func Infof(format string, args ...any) {
	Log.consoleLogger.Infof(format, args...)
	Log.fileLogger.Infof(format, args...)
}

func Warnf(format string, args ...any) {
	Log.consoleLogger.Warnf(format, args...)
	Log.fileLogger.Warnf(format, args...)
}

func Errorf(format string, args ...any) {
	Log.consoleLogger.Errorf(format, args...)
	Log.fileLogger.Errorf(format, args...)
}
