package logger

import (
	"backend/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Interface
type Interface interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type LoggerZap struct {
	*zap.Logger
}

func New(config config.Logger) Interface {
	var level zapcore.Level
	switch config.LogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	encoder := getEncoderLog()

	hook := lumberjack.Logger{
		Filename:   config.FileLogName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&hook),
		),
		level,
	)

	return &LoggerZap{zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)}
}

// format log
func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()

	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// ts -> time
	encodeConfig.TimeKey = "time"
	// from info -> INFO
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}
