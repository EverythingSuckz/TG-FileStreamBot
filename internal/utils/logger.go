package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func InitLogger(debugMode bool) {
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		layout := "02/01/2006 03:04:05 PM"
		if debugMode {
			layout = "02/01/2006 03:04:05.000 PM"
		}
		enc.AppendString(t.Format(layout))
	}
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleConfig.EncodeTime = customTimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)

	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	})

	var consoleLevel zapcore.Level
	if debugMode {
		consoleLevel = zapcore.DebugLevel
	} else {
		consoleLevel = zapcore.InfoLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), consoleLevel),
		zapcore.NewCore(fileEncoder, fileWriter, zapcore.DebugLevel),
	)

	Logger = zap.New(core, zap.AddStacktrace(zapcore.FatalLevel))
}
