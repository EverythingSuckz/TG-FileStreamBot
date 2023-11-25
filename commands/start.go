package commands

import (
	"github.com/celestix/gotgproto/dispatcher"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/ext"
)

func (m *command) LoadStart(dispatcher dispatcher.Dispatcher) {
	defer m.log.Sugar().Info("Loaded start command")
	dispatcher.AddHandler(handlers.NewCommand("start", start))
}

func start(ctx *ext.Context, u *ext.Update) error {
	ctx.Reply(u, "Hi", nil)
	return dispatcher.EndGroups
}
