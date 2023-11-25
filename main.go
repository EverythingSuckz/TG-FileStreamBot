package main

import (
	"EverythingSuckz/fsb/bot"
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/routes"
	"EverythingSuckz/fsb/types"
	"EverythingSuckz/fsb/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const versionString = "v0.0.1"

var startTime time.Time = time.Now()

func main() {
	log := initLogger()
	log.Info("Starting server...")
	config.Load(log)
	router := getRouter(log)

	_, err := bot.StartClient(log)
	if err != nil {
		log.Info(err.Error())
		return
	}
	log.Info("Server started", zap.Int("port", config.ValueOf.Port))
	log.Info("File Stream Bot", zap.String("version", versionString))
	err = router.Run(fmt.Sprintf(":%d", config.ValueOf.Port))
	if err != nil {
		log.Sugar().Fatalln(err)
	}

}

func initLogger() *zap.Logger {
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("02/01/2006 03:04 PM"))
	}
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleConfig.EncodeTime = customTimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	defaultLogLevel := zapcore.DebugLevel

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

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), defaultLogLevel),
		zapcore.NewCore(fileEncoder, fileWriter, defaultLogLevel),
	)

	logger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
	return logger

}

func getRouter(log *zap.Logger) *gin.Engine {
	if config.ValueOf.Dev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(gin.ErrorLogger())
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, types.RootResponse{
			Message: "Server is running.",
			Ok:      true,
			Uptime:  utils.TimeFormat(uint64(time.Since(startTime).Seconds())),
			Version: versionString,
		})
	})
	routes.Load(log, router)
	return router
}
