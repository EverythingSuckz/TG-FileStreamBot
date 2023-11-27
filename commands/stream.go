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
	defer m.log.Sugar().Info("Loaded stream command")
	dispatcher.AddHandler(
		handlers.NewMessage(filters.Message.Media, sendLink),
	)
}

func sendLink(ctx *ext.Context, u *ext.Update) error {
	if u.EffectiveMessage.Media == nil {
		return dispatcher.EndGroups
	}
	if len(u.Entities.Chats) != 0 {
		return dispatcher.EndGroups
	}
	chatId := u.EffectiveChat().GetID()
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
	link := fmt.Sprintf("%s/stream/%d", config.ValueOf.Host, messageID)
	ctx.Reply(u, link, nil)
	return dispatcher.EndGroups
}
