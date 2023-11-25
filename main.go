package main

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/routes"
	"EverythingSuckz/fsb/types"
	"EverythingSuckz/fsb/utils"
	"EverythingSuckz/fsb/bot"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const versionString = "v0.0.1"

var startTime time.Time = time.Now()

func main() {
	log.Println("Starting server...")
	config.Load()
	router := getRouter()

	_, err := bot.StartClient()
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Server started at: http://localhost:%d\n", config.ValueOf.Port)
	err = router.Run(fmt.Sprintf(":%d", config.ValueOf.Port))
	if err != nil {
		log.Println(err)
	}

}

// TODO: Use zap logger
func initLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	return logger
}

func getRouter() *gin.Engine {
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
	routes.Load(router)
	return router
}
