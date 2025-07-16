package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

func Init() {
	logFile := "logs/app.log"
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		panic("Failed to create logs directory: " + err.Error())
	}

	// Encoder config for plain text
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "datetime",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "class", // Rename for clarity
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Console: colored, same encoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// File: plain text
	fileEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	// File writer
	fileWriter, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zap.InfoLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), zap.InfoLevel),
	)

	zapLog = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// L returns a logger enriched with request ID from context
func L(ctx context.Context) *zap.Logger {
	if requestID, ok := ctx.Value("requestID").(string); ok {
		return zapLog.With(zap.String("request_id", requestID))
	}
	return zapLog
}

// Sync flushes any buffered log entries
func Sync() {
	_ = zapLog.Sync()
}
