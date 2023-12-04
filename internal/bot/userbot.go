package bot

import (
	"EverythingSuckz/fsb/config"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
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
		gotgproto.ClientType{
			Phone: "",
		},
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

}
