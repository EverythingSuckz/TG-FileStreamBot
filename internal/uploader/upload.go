package uploader

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/bot"
	"EverythingSuckz/fsb/internal/utils"
	"EverythingSuckz/fsb/pkg/drive115"
	"fmt"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
)

import "go.uber.org/zap"

func Load(log *zap.Logger, dispatcher dispatcher.Dispatcher) {
	log = log.Named("uploader")
	defer log.Info("Initialized uploader command handlers")
	dispatcher.AddHandler(
		handlers.NewCommand("upload", upload),
	)
}

func upload(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveMessage.ReplyToMessage == nil {
		ctx.Reply(u, "Reply to a message to upload.", nil)
		return dispatcher.EndGroups
	}

	driveClient := drive115.NewClient(config.ValueOf.Drive115Cookie)
	uploadInfo, err := driveClient.GetUploadInfo()
	if err != nil {
		ctx.Reply(u, fmt.Sprintf("Failed to get upload info: %s", err), nil)
		return dispatcher.EndGroups
	}

	worker := bot.GetNextWorker()
	file, err := utils.FileFromMessage(ctx, worker.Client, u.EffectiveMessage.ReplyToMessage.GetID())
	if err != nil {
		ctx.Reply(u, fmt.Sprintf("Failed to get file from message: %s", err), nil)
		return dispatcher.EndGroups
	}

	reader, err := utils.NewTelegramReader(ctx, worker.Client, file.Location, 0, file.FileSize-1, file.FileSize)
	if err != nil {
		ctx.Reply(u, fmt.Sprintf("Failed to create telegram reader: %s", err), nil)
		return dispatcher.EndGroups
	}

	result, err := driveClient.Upload(reader, uploadInfo.UploadURL)
	if err != nil {
		ctx.Reply(u, fmt.Sprintf("Failed to upload file: %s", err), nil)
		return dispatcher.EndGroups
	}

	ctx.Reply(u, fmt.Sprintf("File uploaded successfully. File ID: %s", result.FileID), nil)
	return dispatcher.EndGroups
}
