package routes

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (e *allRoutes) LoadGenerate(r *Route) {
	genLog := e.log.Named("Generate")
	defer genLog.Info("Loaded generate route")
	r.Engine.GET("/generate/:messageID", getGenerateRoute)
}

// writeJSON writes a JSON response without HTML-escaping special chars like &.
func writeJSON(ctx *gin.Context, status int, v any) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.Encode(v)
	ctx.Data(status, "application/json; charset=utf-8", buf.Bytes())
}

func getGenerateRoute(ctx *gin.Context) {
	messageIDParam := ctx.Param("messageID")
	messageID, err := strconv.Atoi(messageIDParam)
	if err != nil {
		writeJSON(ctx, http.StatusBadRequest, gin.H{"ok": false, "error": "invalid message ID"})
		return
	}

	token := ctx.Query("token")
	if token == "" {
		token = ctx.GetHeader("X-Bot-Token")
	}
	if token == "" || token != config.ValueOf.BotToken {
		writeJSON(ctx, http.StatusUnauthorized, gin.H{"ok": false, "error": "unauthorized"})
		return
	}

	expParam := ctx.Query("exp")
	var expiresAt int64
	var expiryLabel string

	if expParam == "0" {
		expiresAt = 0
		expiryLabel = "never"
	} else {
		var expiryDuration time.Duration
		if expParam != "" {
			d, err := time.ParseDuration(expParam)
			if err != nil {
				writeJSON(ctx, http.StatusBadRequest, gin.H{"ok": false, "error": fmt.Sprintf("invalid exp duration: %s", expParam)})
				return
			}
			expiryDuration = d
		} else {
			expiryDuration = 24 * time.Hour
		}
		expiresAt = time.Now().Add(expiryDuration).Unix()
		expiryLabel = expiryDuration.String()
	}

	sig := utils.SignURL(messageID, expiresAt)
	link := fmt.Sprintf("%s/stream/%d?exp=%s&sig=%s",
		config.ValueOf.Host,
		messageID,
		strconv.FormatInt(expiresAt, 10),
		sig,
	)

	if ctx.Query("redirect") == "true" {
		ctx.Redirect(http.StatusFound, link)
		return
	}

	writeJSON(ctx, http.StatusOK, gin.H{
		"ok":         true,
		"url":        link,
		"expires_at": expiresAt,
		"expires_in": expiryLabel,
	})
}
