package utils

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/types"
	"context"
	"errors"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func GetTGMessage(ctx context.Context, client *telegram.Client, messageID int) (*tg.Message, error) {
	inputMessageID := tg.InputMessageClass(&tg.InputMessageID{ID: messageID})
	channel, err := GetChannelById(ctx, client)
	if err != nil {
		return nil, err
	}
	messageRequest := tg.ChannelsGetMessagesRequest{Channel: channel, ID: []tg.InputMessageClass{inputMessageID}}
	res, err := client.API().ChannelsGetMessages(ctx, &messageRequest)
	if err != nil {
		return nil, err
	}
	messages := res.(*tg.MessagesChannelMessages)
	message := messages.Messages[0]
	if _, ok := message.(*tg.Message); ok {
		return message.(*tg.Message), nil
	} else {
		return nil, fmt.Errorf("this file was deleted")
	}
}

func GetMessageMedia(ctx context.Context, client *telegram.Client, messageID int) (tg.MessageMediaClass, error) {
	message, err := GetTGMessage(ctx, client, messageID)
	if err != nil {
		return nil, err
	}
	return message.Media, nil
}

func FileFromMedia(media tg.MessageMediaClass) (*types.File, error) {
	switch media := media.(type) {
	case *tg.MessageMediaDocument:
		document, ok := media.Document.AsNotEmpty()
		if !ok {
			return nil, fmt.Errorf("unexpected type %T", media)
		}
		var fileName string
		for _, attribute := range document.Attributes {
			if name, ok := attribute.(*tg.DocumentAttributeFilename); ok {
				fileName = name.FileName
				break
			}
		}
		return &types.File{
			Location: document.AsInputDocumentFileLocation(),
			FileSize: document.Size,
			FileName: fileName,
			MimeType: document.MimeType,
			ID:       document.ID,
		}, nil
		// TODO: add photo support
	}
	return nil, fmt.Errorf("unexpected type %T", media)
}

func GetChannelById(ctx context.Context, client *telegram.Client) (*tg.InputChannel, error) {
	channel := &tg.InputChannel{}
	inputChannel := &tg.InputChannel{
		ChannelID: config.ValueOf.LogChannelID,
	}
	channels, err := client.API().ChannelsGetChannels(ctx, []tg.InputChannelClass{inputChannel})
	if err != nil {
		return nil, err
	}
	if len(channels.GetChats()) == 0 {
		return nil, errors.New("no channels found")
	}
	channel = channels.GetChats()[0].(*tg.Channel).AsInput()
	return channel, nil
}
