package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
	Log  *zap.Logger
)

// Init initializes the global Zap logger.
// Panics if called more than once (use Init in main).
func Init(development bool) {
	once.Do(func() {
		var cfg zap.Config
		if development {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		} else {
			cfg = zap.NewProductionConfig()
			cfg.EncoderConfig.TimeKey = "timestamp"
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		var err error
		Log, err = cfg.Build(zap.AddCallerSkip(0))
		if err != nil {
			panic("failed to initialize zap logger: " + err.Error())
		}
	})
}

// FromEnv initializes the logger based on APP_ENV environment variable.
// APP_ENV=production → structured JSON; anything else → development (colored console).
func FromEnv() {
	dev := os.Getenv("APP_ENV") != "production"
	Init(dev)
}

// Named returns a named child logger of the global logger.
// Panics if the global logger is not initialized yet.
func Named(name string) *zap.Logger {
	if Log == nil {
		panic("logger not initialized: call logger.FromEnv() or logger.Init() first")
	}
	return Log.Named(name)
}
