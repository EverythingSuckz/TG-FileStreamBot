package bot

import (
	"EverythingSuckz/fsb/config"
	"errors"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type UserBotStruct struct {
	log    *zap.Logger
	client *gotgproto.Client
}

var UserBot *UserBotStruct = &UserBotStruct{}

func StartUserBot(l *zap.Logger) {
	log := l.Named("USERBOT")
	if config.ValueOf.UserSession == "" {
		log.Warn("User session is empty")
		return
	}
	log.Sugar().Infoln("Starting userbot")
	client, err := gotgproto.NewClient(
		int(config.ValueOf.ApiID),
		config.ValueOf.ApiHash,
		gotgproto.ClientTypePhone(""),
		&gotgproto.ClientOpts{
			Session:          sessionMaker.PyrogramSession(config.ValueOf.UserSession),
			DisableCopyright: true,
		},
	)
	if err != nil {
		log.Error("Failed to start userbot", zap.Error(err))
		return
	}
	UserBot.log = log
	UserBot.client = client
	log.Info("Userbot started", zap.String("username", client.Self.Username), zap.String("FirstName", client.Self.FirstName), zap.String("LastName", client.Self.LastName))
	if err := UserBot.AddBotsAsAdmins(); err != nil {
		log.Error("Failed to add bots as admins", zap.Error(err))
		return
	}
}

func (u *UserBotStruct) AddBotsAsAdmins() error {
	u.log.Info("Preparing to add bots as admins")
	ctx := u.client.CreateContext()
	channel := config.ValueOf.LogChannelID
	channelInfos, err := u.client.API().ChannelsGetChannels(
		ctx,
		[]tg.InputChannelClass{
			&tg.InputChannel{
				ChannelID: channel,
			},
		},
	)
	if err != nil {
		u.log.Error("Failed to get channel info", zap.Error(err))
		return errors.New("failed to get channel info")
	}
	if len(channelInfos.GetChats()) == 0 {
		return errors.New("no channels found")
	}
	inputChannel := channelInfos.GetChats()[0].(*tg.Channel).AsInput()
	currentAdmins := []int64{}
	admins, err := u.client.API().ChannelsGetParticipants(ctx, &tg.ChannelsGetParticipantsRequest{
		Channel: inputChannel,
		Filter:  &tg.ChannelParticipantsAdmins{},
		Offset:  0,
		Limit:   100,
	})
	if err != nil {
		u.log.Error("Failed to get admins", zap.Error(err))
		return err
	}
	for _, admin := range admins.(*tg.ChannelsChannelParticipants).Participants {
		if user, ok := admin.(*tg.ChannelParticipantAdmin); ok {
			currentAdmins = append(currentAdmins, user.UserID)
		}
	}
	for _, bot := range Workers.Bots {
		isAdmin := false
		for _, admin := range currentAdmins {
			if admin == bot.Self.ID {
				u.log.Sugar().Infof("Bot @%s is already an admin", bot.Self.Username)
				isAdmin = true
				continue
			}
		}
		if isAdmin {
			continue
		}
		botInfo, err := ctx.ResolveUsername(bot.Self.Username)
		if err != nil {
			u.log.Warn(err.Error())
		}
		_, err = u.client.API().ChannelsEditAdmin(
			u.client.CreateContext().Context,
			&tg.ChannelsEditAdminRequest{
				Channel: inputChannel,
				UserID:  botInfo.GetInputUser(),
				AdminRights: tg.ChatAdminRights{
					PostMessages: true,
				},
				Rank: "admin",
			},
		)
		if err != nil {
			u.log.Sugar().Warnf("Failed to add @%s as admin", bot.Self.Username)
			u.log.Warn(err.Error())
		}
		u.log.Sugar().Infof("Added @%s as admin", bot.Self.Username)
	}
	return nil
}
