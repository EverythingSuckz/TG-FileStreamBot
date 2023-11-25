package bot

import (
	"EverythingSuckz/fsb/commands"
	"EverythingSuckz/fsb/config"
	"log"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/sessionMaker"
)

func StartClient() (*gotgproto.Client, error) {
	client, err := gotgproto.NewClient(
		int(config.ValueOf.ApiID),
		config.ValueOf.ApiHash,
		gotgproto.ClientType{
			BotToken: config.ValueOf.BotToken,
		},
		&gotgproto.ClientOpts{
			Session:          sessionMaker.NewSession("fsb", sessionMaker.Session),
			DisableCopyright: true,
		},
	)
	if err != nil {
		return nil, err
	}
	commands.Load(client.Dispatcher)
	log.Println("Client started")
	return client, nil
}
