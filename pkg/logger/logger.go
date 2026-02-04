package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// Log is the package-level zerolog logger configured by Init or package init.
var Log zerolog.Logger

func init() {
	// provide a sensible default so other packages can log before explicit Init()
	zerolog.TimeFieldFormat = time.RFC3339
	w := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	l := zerolog.New(w).With().Timestamp().Logger().Level(zerolog.InfoLevel)
	Log = l
	zlog.Logger = l
}

// Init configures the global logger. It reads LOG_LEVEL from the environment
// and overrides the package logger. Call this early in `main` when possible.
func Init() {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if level == "" {
		level = "info"
	}

	zerolog.TimeFieldFormat = time.RFC3339
	w := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	l := zerolog.New(w).With().Timestamp().Logger().Level(parseLevel(level))

	Log = l
	zlog.Logger = l
}

func parseLevel(s string) zerolog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return zerolog.DebugLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}
