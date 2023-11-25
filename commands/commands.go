package commands

import (
	"log"
	"reflect"

	"github.com/celestix/gotgproto/dispatcher"
)

type command struct {
}

func Load(dispatcher dispatcher.Dispatcher) {
	defer log.Println("Initialized")
	Type := reflect.TypeOf(&command{})
	Value := reflect.ValueOf(&command{})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(dispatcher)})
	}
}
