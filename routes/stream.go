package routes

import (
	"EverythingSuckz/fsb/bot"
	"EverythingSuckz/fsb/utils"
	"fmt"
	"io"
	"net/http"
	"strconv"

	range_parser "github.com/quantumsheep/range-parser"

	"github.com/gin-gonic/gin"
)

func (e *allRoutes) LoadHome(r *Route) {
	defer e.log.Info("Loaded stream route")
	r.Engine.GET("/stream/:messageID", getStreamRoute)
}

func getStreamRoute(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request

	messageIDParm := ctx.Param("messageID")
	messageID, err := strconv.Atoi(messageIDParm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.Header("Accept-Ranges", "bytes")
	var start, end int64
	rangeHeader := r.Header.Get("Range")

	messageMedia, err := utils.GetMessageMedia(ctx, bot.Bot.Client, messageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, err := utils.FileFromMedia(messageMedia)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if rangeHeader == "" {
		start = 0
		end = file.FileSize - 1
		w.WriteHeader(http.StatusOK)
	} else {
		ranges, err := range_parser.Parse(file.FileSize, r.Header.Get("Range"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		start = ranges[0].Start
		end = ranges[0].End
		ctx.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, file.FileSize))
		w.WriteHeader(http.StatusPartialContent)
	}

	contentLength := end - start + 1
	mimeType := file.MimeType

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	ctx.Header("Content-Type", mimeType)
	ctx.Header("Content-Length", strconv.FormatInt(contentLength, 10))

	disposition := "inline"

	if ctx.Query("d") == "true" {
		disposition = "attachment"
	}

	ctx.Header("Content-Disposition", fmt.Sprintf("%s; filename=\"%s\"", disposition, file.FileName))

	if r.Method != "HEAD" {
		parts, err := utils.GetParts(ctx, bot.Bot.Client, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		parts = utils.RangedParts(parts, start, end)
		lr, _ := utils.NewLinearReader(ctx, bot.Bot.Client, parts, contentLength)
		io.CopyN(w, lr, contentLength)
	}
}
