package bot

import (
	"EverythingSuckz/fsb/commands"
	"EverythingSuckz/fsb/config"

	"go.uber.org/zap"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
)

var Bot *gotgproto.Client

func StartClient(log *zap.Logger) (*gotgproto.Client, error) {
	client, err := gotgproto.NewClient(
		int(config.ValueOf.ApiID),
		config.ValueOf.ApiHash,
		gotgproto.ClientType{
			BotToken: config.ValueOf.BotToken,
		},
		&gotgproto.ClientOpts{
			Session:          sessionMaker.SqliteSession("fsb"),
			DisableCopyright: true,
		},
	)
	if err != nil {
		return nil, err
	}
	commands.Load(log, client.Dispatcher)
	log.Info("Client started", zap.String("username", client.Self.Username))
	Bot = client
	return client, nil
}
