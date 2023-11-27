package commands

import (
	"reflect"

	"github.com/celestix/gotgproto/dispatcher"
	"go.uber.org/zap"
)

type command struct {
	log *zap.Logger
}

func Load(log *zap.Logger, dispatcher dispatcher.Dispatcher) {
	log = log.Named("commands")
	defer log.Info("Initialized all command handlers")
	Type := reflect.TypeOf(&command{log})
	Value := reflect.ValueOf(&command{log})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(dispatcher)})
	}
}
