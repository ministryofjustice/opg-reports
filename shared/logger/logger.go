// Package logger provides a wrapper for slog configuration
//
// It looks for env variables (LOG_LEVEL, LOG_AS, LOG_TO) to determine level, type and location
// the logs should go to. By default, it assumed "info", "txt" and "stdout" respectively.
//
// Call `LogSetup()` to configure the default slog logger
package logger

import (
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
)

var (
	logLevels = map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	logAsChoices    = []string{"text", "json"}
	logToChoices    = []string{"stdout", "file"}
	logFile         *os.File
	currentlogLevel slog.Level
)

func Level() slog.Level {
	return currentlogLevel
}

// LogSetup handles config for the default logger, checking env vars to determine log levels,
// log type (json or text) and log destination (stdout / file)
func LogSetup() {
	var (
		level          string = os.Getenv("LOG_LEVEL") // "info"
		as             string = os.Getenv("LOG_AS")    // "text"
		to             string = os.Getenv("LOG_TO")    // "stdout"
		validAsChoice  bool
		validToChoice  bool
		out            io.Writer  = os.Stdout
		logLevel       slog.Level = slog.LevelError
		handlerOptions *slog.HandlerOptions
		log            *slog.Logger
	)

	// setup log level
	level = strings.ToLower(level)
	if l, ok := logLevels[level]; ok {
		logLevel = l
	} else {
		logLevel = logLevels["info"]
	}
	currentlogLevel = logLevel
	// setup log as
	validAsChoice = slices.Contains(logAsChoices, as)
	if !validAsChoice {
		as = "text"
	}
	// setup to
	validToChoice = slices.Contains(logToChoices, to)
	if !validToChoice {
		to = "stdout"
	}

	handlerOptions = &slog.HandlerOptions{AddSource: false, Level: logLevel}
	// if chosen to change output to file, open the file and adjust out
	if validToChoice && to == "file" {
		logFile, _ = os.OpenFile("log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		out = logFile
	}

	if validAsChoice && as == "json" {
		log = slog.New(slog.NewJSONHandler(out, handlerOptions))
	} else {
		log = slog.New(slog.NewTextHandler(out, handlerOptions))
	}
	slog.SetDefault(log)

}
