package routes

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/bot"
	"EverythingSuckz/fsb/internal/utils"
	"EverythingSuckz/fsb/pkg/drive115"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (e *allRoutes) LoadUpload(r *Route) {
	log := e.log.Named("Upload")
	defer log.Info("Loaded upload route")
	r.Engine.GET("/upload/:messageID", getUploadRoute)
}

func getUploadRoute(ctx *gin.Context) {
	w := ctx.Writer

	messageIDParm := ctx.Param("messageID")
	messageID, err := strconv.Atoi(messageIDParm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authHash := ctx.Query("hash")
	if authHash == "" {
		http.Error(w, "missing hash param", http.StatusBadRequest)
		return
	}

	worker := bot.GetNextWorker()

	file, err := utils.FileFromMessage(ctx, worker.Client, messageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expectedHash := utils.PackFile(
		file.FileName,
		file.FileSize,
		file.MimeType,
		file.ID,
	)
	if !utils.CheckHash(authHash, expectedHash) {
		http.Error(w, "invalid hash", http.StatusBadRequest)
		return
	}

	driveClient := drive115.NewClient(config.ValueOf.Drive115Cookie)
	uploadInfo, err := driveClient.GetUploadInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get upload info: %s", err), http.StatusInternalServerError)
		return
	}

	reader, err := utils.NewTelegramReader(ctx, worker.Client, file.Location, 0, file.FileSize-1, file.FileSize)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create telegram reader: %s", err), http.StatusInternalServerError)
		return
	}

	result, err := driveClient.Upload(reader, uploadInfo.UploadURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upload file: %s", err), http.StatusInternalServerError)
		return
	}

	ctx.String(http.StatusOK, "File uploaded successfully. File ID: %s", result.FileID)
}
