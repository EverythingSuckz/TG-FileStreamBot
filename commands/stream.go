package commands

import (
	"fmt"

	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/utils"

	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/storage"
	"github.com/gotd/td/tg"
)

func (m *command) LoadStream(dispatcher dispatcher.Dispatcher) {
	log := m.log.Named("start")
	defer log.Sugar().Info("Loaded")
	dispatcher.AddHandler(
		handlers.NewMessage(filters.Message.Media, sendLink),
	)
}

func sendLink(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveMessage.Media == nil {
		return dispatcher.EndGroups
	}
	chatId := u.EffectiveChat().GetID()
	peerChatId := storage.GetPeerById(chatId)
	if peerChatId.Type != int(storage.TypeUser) {
		return dispatcher.EndGroups
	}
	peer := storage.GetPeerById(config.ValueOf.LogChannelID)
	switch storage.EntityType(peer.Type) {
	case storage.TypeChat:
		return dispatcher.EndGroups
	case storage.TypeUser:
		return dispatcher.EndGroups
	}
	update, err := ctx.ForwardMessages(
		chatId,
		config.ValueOf.LogChannelID,
		&tg.MessagesForwardMessagesRequest{
			FromPeer: &tg.InputPeerChat{ChatID: chatId},
			ID:       []int{u.EffectiveMessage.ID},
			ToPeer:   &tg.InputPeerChannel{ChannelID: config.ValueOf.LogChannelID, AccessHash: peer.AccessHash},
		},
	)
	if err != nil {
		utils.Logger.Sugar().Error(err)
		ctx.Reply(u, fmt.Sprintf("Error - %s", err.Error()), nil)
		return dispatcher.EndGroups
	}
	messageID := update.(*tg.Updates).Updates[0].(*tg.UpdateMessageID).ID
	if err != nil {
		utils.Logger.Sugar().Error(err)
		ctx.Reply(u, fmt.Sprintf("Error - %s", err.Error()), nil)
		return dispatcher.EndGroups
	}
	doc := update.(*tg.Updates).Updates[1].(*tg.UpdateNewChannelMessage).Message.(*tg.Message).Media
	file, err := utils.FileFromMedia(doc)
	if err != nil {
		ctx.Reply(u, fmt.Sprintf("Error - %s", err.Error()), nil)
		return dispatcher.EndGroups
	}
	fullHash := utils.PackFile(
		file.FileName,
		file.FileSize,
		file.MimeType,
		file.ID,
	)

	hash := utils.GetShortHash(fullHash)
	link := fmt.Sprintf("%s/stream/%d?hash=%s", config.ValueOf.Host, messageID, hash)
	ctx.Reply(u, link, nil)
	return dispatcher.EndGroups
}
