package utils

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/cache"
	"EverythingSuckz/fsb/internal/types"
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

// https://stackoverflow.com/a/70802740/15807350
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func GetTGMessage(ctx context.Context, client *gotgproto.Client, messageID int) (*tg.Message, error) {
	inputMessageID := tg.InputMessageClass(&tg.InputMessageID{ID: messageID})
	channel, err := GetLogChannelPeer(ctx, client.API(), client.PeerStorage)
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

func FileFromMessage(ctx context.Context, client *gotgproto.Client, messageID int) (*types.File, error) {
	key := fmt.Sprintf("file:%d:%d", messageID, client.Self.ID)
	log := Logger.Named("GetMessageMedia")
	var cachedMedia types.File
	err := cache.GetCache().Get(key, &cachedMedia)
	if err == nil {
		log.Debug("Using cached media message properties", zap.Int("messageID", messageID), zap.Int64("clientID", client.Self.ID))
		return &cachedMedia, nil
	}
	log.Debug("Fetching file properties from message ID", zap.Int("messageID", messageID), zap.Int64("clientID", client.Self.ID))
	message, err := GetTGMessage(ctx, client, messageID)
	if err != nil {
		return nil, err
	}
	file, err := FileFromMedia(message.Media)
	if err != nil {
		return nil, err
	}
	err = cache.GetCache().Set(
		key,
		file,
		3600,
	)
	if err != nil {
		return nil, err
	}
	return file, nil
	// TODO: add photo support
}

func GetLogChannelPeer(ctx context.Context, api *tg.Client, peerStorage *storage.PeerStorage) (*tg.InputChannel, error) {
	cachedInputPeer := peerStorage.GetInputPeerById(config.ValueOf.LogChannelID)

	switch peer := cachedInputPeer.(type) {
	case *tg.InputPeerEmpty:
		break
	case *tg.InputPeerChannel:
		return &tg.InputChannel{
			ChannelID:  peer.ChannelID,
			AccessHash: peer.AccessHash,
		}, nil
	default:
		return nil, errors.New("unexpected type of input peer")
	}
	inputChannel := &tg.InputChannel{
		ChannelID: config.ValueOf.LogChannelID,
	}
	channels, err := api.ChannelsGetChannels(ctx, []tg.InputChannelClass{inputChannel})
	if err != nil {
		return nil, err
	}
	if len(channels.GetChats()) == 0 {
		return nil, errors.New("no channels found")
	}
	channel, ok := channels.GetChats()[0].(*tg.Channel)
	if !ok {
		return nil, errors.New("type assertion to *tg.Channel failed")
	}
	// Bruh, I literally have to call library internal functions at this point
	peerStorage.AddPeer(channel.GetID(), channel.AccessHash, storage.TypeChannel, "")
	return channel.AsInput(), nil
}

func ForwardMessages(ctx *ext.Context, fromChatId, toChatId int64, messageID int) (*tg.Updates, error) {
	fromPeer := ctx.PeerStorage.GetInputPeerById(fromChatId)
	if fromPeer.Zero() {
		return nil, fmt.Errorf("fromChatId: %d is not a valid peer", fromChatId)
	}
	toPeer, err := GetLogChannelPeer(ctx, ctx.Raw, ctx.PeerStorage)
	if err != nil {
		return nil, err
	}
	update, err := ctx.Raw.MessagesForwardMessages(ctx, &tg.MessagesForwardMessagesRequest{
		RandomID: []int64{rand.Int63()},
		FromPeer: fromPeer,
		ID:       []int{messageID},
		ToPeer:   &tg.InputPeerChannel{ChannelID: toPeer.ChannelID, AccessHash: toPeer.AccessHash},
	})
	if err != nil {
		return nil, err
	}
	return update.(*tg.Updates), nil
}
