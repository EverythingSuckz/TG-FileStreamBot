package main

import (
	"EverythingSuckz/fsb/bot"
	"EverythingSuckz/fsb/cache"
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/routes"
	"EverythingSuckz/fsb/types"
	"EverythingSuckz/fsb/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const versionString = "v0.0.0"

var startTime time.Time = time.Now()

func main() {
	utils.InitLogger()
	log := utils.Logger
	mainLogger := log.Named("Main")
	mainLogger.Info("Starting server")
	config.Load(log)
	router := getRouter(log)
	
	_, err := bot.StartClient(log)
	if err != nil {
		log.Info(err.Error())
		return
	}
	cache.InitCache(log)
	mainLogger.Info("Server started", zap.Int("port", config.ValueOf.Port))
	mainLogger.Info("File Stream Bot", zap.String("version", versionString))
	err = router.Run(fmt.Sprintf(":%d", config.ValueOf.Port))
	if err != nil {
		mainLogger.Sugar().Fatalln(err)
	}

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
